package main

import (
	"arq-soft-II/config/rabbitmq"
	"fmt"
	"time"
)

func main() {
	r, err := rabbitmq.New("amqp://admin:admin@localhost:5672/")
	if err != nil {
		panic(err)
	}
	defer r.Close()

	fmt.Println("✅ Conexión exitosa con RabbitMQ!")

	err = r.DeclareSetup("entity.events", "search-sync", "activities.*")
	if err != nil {
		panic(err)
	}

	msg := []byte(`{"id":1, "action":"created"}`)
	if err := r.Publish("entity.events", "activities.created", msg); err != nil {
		panic(err)
	}

	fmt.Println("📨 Mensaje enviado correctamente!")
	time.Sleep(2 * time.Second)
}
