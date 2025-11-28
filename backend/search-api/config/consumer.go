package config

import (
	"encoding/json"
	"log"

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
		"entity.events", // name
		"topic",         // type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		log.Fatal("笶・Error declarando exchange:", err)
	}

	// Declarar queue
	q, err := ch.QueueDeclare(
		"search-sync", // name
		true,          // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		log.Fatal("笶・Error declarando queue:", err)
	}

	// Bind queue to exchange
	err = ch.QueueBind(
		q.Name,          // queue name
		"activities.*",  // routing key
		"entity.events", // exchange
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		log.Fatal("笶・Error binding queue:", err)
	}

	// Consume mensajes
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatal("笶・Error al consumir mensajes:", err)
	}

	log.Println("笨・RabbitMQ Consumer iniciado - Esperando eventos de actividades...")

	go func() {
		for msg := range msgs {
			log.Printf("陶 [RabbitMQ] Mensaje recibido: %s", string(msg.Body))

			// Parsear el evento
			var event ActivityEvent
			if err := json.Unmarshal(msg.Body, &event); err != nil {
				log.Printf("笶・Error parseando mensaje: %v", err)
				msg.Nack(false, false) // No reencolar
				continue
			}

			// Procesar segﾃｺn el tipo de acciﾃｳn
			var processErr error
			switch event.Action {
			case "created":
				log.Printf("・ Procesando creaciﾃｳn de actividad: %s", event.ActivityID)
				processErr = searchService.IndexActivity(event.Data)

			case "updated":
				log.Printf("売 Procesando actualizaciﾃｳn de actividad: %s", event.ActivityID)
				processErr = searchService.UpdateActivity(event.Data)

			case "deleted":
				log.Printf("卵・・ Procesando eliminaciﾃｳn de actividad: %s", event.ActivityID)
				processErr = searchService.DeleteActivity(event.ActivityID)

			default:
				log.Printf("笞・・ Acciﾃｳn desconocida: %s", event.Action)
			}

			if processErr != nil {
				log.Printf("笶・Error procesando evento: %v", processErr)
				msg.Nack(false, true) // Reencolar para reintentar
			} else {
				log.Printf("笨・Evento procesado correctamente: %s - %s", event.Action, event.ActivityID)
				msg.Ack(false)
			}
		}
	}()
}
