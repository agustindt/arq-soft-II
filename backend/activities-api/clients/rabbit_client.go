package clients

import (
	"arq-soft-II/backend/activities-api/domain"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
)

const (
	encodingJSON = "application/json"
	encodingUTF8 = "UTF-8"
)

type RabbitMQClient struct {
	connection *amqp091.Connection
	channel    *amqp091.Channel
	exchange   string
	repository ActivityRepository
}

// ActivityRepository interface para obtener la actividad completa
type ActivityRepository interface {
	GetByID(ctx context.Context, id string) (domain.Activity, error)
}

// NewRabbitMQClient crea un nuevo cliente de RabbitMQ
func NewRabbitMQClient(user, password, host, port, exchange string, repo ActivityRepository) *RabbitMQClient {
	connStr := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, password, host, port)
	connection, err := amqp091.Dial(connStr)
	if err != nil {
		log.Fatalf("❌ Failed to connect to RabbitMQ: %v", err)
	}

	channel, err := connection.Channel()
	if err != nil {
		log.Fatalf("❌ Failed to open a channel: %v", err)
	}

	// Declarar exchange tipo topic
	err = channel.ExchangeDeclare(
		exchange, // name
		"topic",  // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		log.Fatalf("❌ Failed to declare exchange: %v", err)
	}

	log.Printf("✅ Connected to RabbitMQ - Exchange: %s", exchange)

	return &RabbitMQClient{
		connection: connection,
		channel:    channel,
		exchange:   exchange,
		repository: repo,
	}
}

// Publish publica un evento de actividad a RabbitMQ
func (r *RabbitMQClient) Publish(ctx context.Context, action string, activityID string) error {
	// Obtener los datos completos de la actividad
	var activityData domain.Activity
	var err error

	if action != "deleted" {
		activityData, err = r.repository.GetByID(ctx, activityID)
		if err != nil {
			return fmt.Errorf("error fetching activity data: %w", err)
		}
	}

	// Preparar el mensaje con toda la información
	message := map[string]interface{}{
		"action":      action,
		"activity_id": activityID,
		"timestamp":   time.Now().UTC().Format(time.RFC3339),
	}

	// Si no es delete, incluir todos los datos
	if action != "deleted" {
		message["data"] = map[string]interface{}{
			"id":           activityData.ID,
			"name":         activityData.Name,
			"description":  activityData.Description,
			"category":     activityData.Category,
			"difficulty":   activityData.Difficulty,
			"location":     activityData.Location,
			"price":        activityData.Price,
			"duration":     activityData.Duration,
			"max_capacity": activityData.MaxCapacity,
			"instructor":   activityData.Instructor,
			"schedule":     activityData.Schedule,
			"equipment":    activityData.Equipment,
			"image_url":    activityData.ImageURL,
			"is_active":    activityData.IsActive,
			"created_by":   activityData.CreatedBy,
			"created_at":   activityData.CreatedAt.Format(time.RFC3339),
			"updated_at":   activityData.UpdatedAt.Format(time.RFC3339),
		}
	}

	bytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error marshalling message to JSON: %w", err)
	}

	// Routing key: activities.{action}
	routingKey := fmt.Sprintf("activities.%s", action)

	err = r.channel.PublishWithContext(
		ctx,
		r.exchange, // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp091.Publishing{
			ContentType:     encodingJSON,
			ContentEncoding: encodingUTF8,
			DeliveryMode:    amqp091.Persistent,
			MessageId:       uuid.New().String(),
			Timestamp:       time.Now().UTC(),
			AppId:           "activities-api",
			Body:            bytes,
		},
	)

	if err != nil {
		return fmt.Errorf("error publishing message to RabbitMQ: %w", err)
	}

	return nil
}

// Close cierra el channel y la conexión a RabbitMQ
func (r *RabbitMQClient) Close() error {
	var errRet error
	if r.channel != nil {
		if err := r.channel.Close(); err != nil {
			log.Printf("⚠️  Error closing rabbit channel: %v", err)
			errRet = err
		}
	}
	if r.connection != nil {
		if err := r.connection.Close(); err != nil {
			log.Printf("⚠️  Error closing rabbit connection: %v", err)
			errRet = err
		}
	}
	return errRet
}
