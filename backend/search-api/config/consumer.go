package config

import (
	"log"

	"arq-soft-II/backend/search-api/services"
	"arq-soft-II/config/rabbitmq"
)

func StartRabbitConsumer(mq *rabbitmq.Rabbit, searchService *services.SearchService) {
	msgs, err := mq.Consume("search-sync")
	if err != nil {
		log.Fatal("❌ Error al consumir mensajes de RabbitMQ:", err)
	}

	go func() {
		for msg := range msgs {
			log.Println("📩 [RabbitMQ] Mensaje recibido:", string(msg.Body))
			searchService.Search("reindex") // simulamos actualización del índice
			msg.Ack(false)
		}
	}()
}
