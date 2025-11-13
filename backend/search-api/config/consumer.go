package config

import (
	"encoding/json"
	"log"

	"arq-soft-II/backend/search-api/domain"
	"arq-soft-II/backend/search-api/services"
	"arq-soft-II/config/rabbitmq"
)

// ActivityEvent representa el evento recibido desde Activities API
type ActivityEvent struct {
	Action     string          `json:"action"`
	ActivityID string          `json:"activity_id"`
	Timestamp  string          `json:"timestamp"`
	Data       domain.Activity `json:"data"`
}

func StartRabbitConsumer(mq *rabbitmq.Rabbit, searchService *services.SearchService) {
	msgs, err := mq.Consume("search-sync")
	if err != nil {
		log.Fatal("❌ Error al consumir mensajes de RabbitMQ:", err)
	}

	log.Println("✅ RabbitMQ Consumer iniciado - Esperando eventos de actividades...")

	go func() {
		for msg := range msgs {
			log.Printf("📩 [RabbitMQ] Mensaje recibido: %s", string(msg.Body))

			// Parsear el evento
			var event ActivityEvent
			if err := json.Unmarshal(msg.Body, &event); err != nil {
				log.Printf("❌ Error parseando mensaje: %v", err)
				msg.Nack(false, false) // No reencolar
				continue
			}

			// Procesar según el tipo de acción
			var processErr error
			switch event.Action {
			case "created":
				log.Printf("🆕 Procesando creación de actividad: %s", event.ActivityID)
				processErr = searchService.IndexActivity(event.Data)

			case "updated":
				log.Printf("🔄 Procesando actualización de actividad: %s", event.ActivityID)
				processErr = searchService.UpdateActivity(event.Data)

			case "deleted":
				log.Printf("🗑️  Procesando eliminación de actividad: %s", event.ActivityID)
				processErr = searchService.DeleteActivity(event.ActivityID)

			default:
				log.Printf("⚠️  Acción desconocida: %s", event.Action)
			}

			if processErr != nil {
				log.Printf("❌ Error procesando evento: %v", processErr)
				msg.Nack(false, true) // Reencolar para reintentar
			} else {
				log.Printf("✅ Evento procesado correctamente: %s - %s", event.Action, event.ActivityID)
				msg.Ack(false)
			}
		}
	}()
}
