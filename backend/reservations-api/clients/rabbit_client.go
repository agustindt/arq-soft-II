package clients

import (
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
	queue      *amqp091.Queue
}

func NewRabbitMQClient(user, password, queueName, host, port string) *RabbitMQClient {
	connStr := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, password, host, port) // 👈 %s
	connection, err := amqp091.Dial(connStr)
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", err) // 👈 %v, no %w
	}
	channel, err := connection.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel: %v", err)
	}
	queue, err := channel.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		log.Fatalf("failed to declare a queue: %v", err)
	}
	return &RabbitMQClient{connection: connection, channel: channel, queue: &queue}
}

func (r RabbitMQClient) Publish(ctx context.Context, action string, reservaID string) error {
	message := map[string]interface{}{
		"action":     action,
		"reserva_id": reservaID,
	}

	bytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error marshalling message to JSON: %w", err)
	}

	if err := r.channel.PublishWithContext(ctx, "", r.queue.Name, false, false, amqp091.Publishing{
		ContentType:     encodingJSON,
		ContentEncoding: encodingUTF8,
		DeliveryMode:    amqp091.Transient,
		MessageId:       uuid.New().String(),
		Timestamp:       time.Now().UTC(),
		AppId:           "reservations-api",
		Body:            bytes,
	}); err != nil {
		return fmt.Errorf("error publishing message to RabbitMQ: %w", err)
	}
	return nil
}

// Close cierra el channel y la conexión a RabbitMQ.
func (r *RabbitMQClient) Close() error {
	var errRet error
	if r.channel != nil {
		if err := r.channel.Close(); err != nil {
			log.Printf("error closing rabbit channel: %v", err)
			errRet = err
		}
	}
	if r.connection != nil {
		if err := r.connection.Close(); err != nil {
			log.Printf("error closing rabbit connection: %v", err)
			errRet = err
		}
	}
	return errRet
}
