package config

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"search-api/clients"
	"search-api/domain"
	"search-api/services"

	amqp "github.com/rabbitmq/amqp091-go"
)

// ActivityEvent representa el evento recibido desde Activities API
type ActivityEvent struct {
	Action     string          `json:"action"`
	ActivityID string          `json:"activity_id"`
	Timestamp  string          `json:"timestamp"`
	Data       domain.Activity `json:"data"`
}

func StartRabbitConsumer(conn *amqp.Connection, searchService *services.SearchService) {
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Error al crear channel:", err)
	}
	defer ch.Close()

	// Declarar exchange
	err = ch.ExchangeDeclare(
		"entity.events",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("❌ Error declaring exchange:", err)
	}

	// Declarar queue
	q, err := ch.QueueDeclare(
		"search-activities",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("❌ Error declaring queue:", err)
	}

	// Binding
	err = ch.QueueBind(
		q.Name,
		"activities.*",
		"entity.events",
		false,
		nil,
	)
	if err != nil {
		log.Fatal("❌ Error binding queue:", err)
	}

	// Consumir mensajes
	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("❌ Error al consumir mensajes:", err)
	}

	log.Println("🚀 RabbitMQ Consumer iniciado - Esperando eventos de actividades...")

	activitiesAPI := os.Getenv("ACTIVITIES_API_URL")
	if activitiesAPI == "" {
		activitiesAPI = "http://activities-api:8082"
	}

	// Retry simple con 3 intentos progresivos
	fetchActivity := func(id string) (domain.Activity, error) {
		var lastErr error
		for attempt := 0; attempt < 3; attempt++ {
			httpCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			activity, err := clients.GetActivityByID(httpCtx, activitiesAPI, id)
			cancel()

			if err == nil {
				return activity, nil
			}

			lastErr = err
			time.Sleep(time.Duration(attempt+1) * 400 * time.Millisecond)
		}
		return domain.Activity{}, lastErr
	}

	go func() {
		for msg := range msgs {
			log.Printf("📩 [RabbitMQ] Mensaje recibido: %s", string(msg.Body))

			var event ActivityEvent
			if err := json.Unmarshal(msg.Body, &event); err != nil {
				log.Printf("❌ Error parseando mensaje: %v", err)
				msg.Nack(false, false)
				continue
			}

			var processErr error

			switch event.Action {
			case "created":
				log.Printf("🟩 Procesando creación de actividad: %s", event.ActivityID)

				activity, err := fetchActivity(event.ActivityID)
				if err != nil {
					processErr = err
					break
				}

				processErr = searchService.IndexActivity(activity)

			case "updated":
				log.Printf("🟦 Procesando actualización de actividad: %s", event.ActivityID)

				activity, err := fetchActivity(event.ActivityID)
				if err != nil {
					processErr = err
					break
				}

				processErr = searchService.UpdateActivity(activity)

			case "deleted":
				log.Printf("🟥 Procesando eliminación de actividad: %s", event.ActivityID)
				processErr = searchService.DeleteActivity(event.ActivityID)

			default:
				log.Printf("⚠️ Acción desconocida: %s", event.Action)
			}

			if processErr != nil {
				log.Printf("❌ Error procesando evento: %v", processErr)
				msg.Nack(false, true)
			} else {
				log.Printf("✅ Evento procesado correctamente: %s - %s", event.Action, event.ActivityID)
				msg.Ack(false)
			}
		}
	}()
}
