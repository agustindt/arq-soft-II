package rabbitmq

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Rabbit estructura principal de conexión
type Rabbit struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

// New crea la conexión principal con RabbitMQ
func New(url string) (*Rabbit, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &Rabbit{Conn: conn, Channel: ch}, nil
}

// "buzón" Se envían los mensajes ahí.
func (r *Rabbit) DeclareSetup(exchange, queue, routingKey string) error {
	if err := r.Channel.ExchangeDeclare(
		exchange, // nombre
		"topic",  // tipo de exchange
		true,     // durable
		false,    // auto-delete
		false,    // internal
		false,    // no-wait
		nil,      // argumentos
	); err != nil {
		return err
	}
	// Es la cola donde se almacenan los mensajes hasta que un consumidor los lea.
	q, err := r.Channel.QueueDeclare(
		queue,
		true, false, false, false, nil,
	)
	if err != nil {
		return err
	}
	// Close cierra la conexión
	return r.Channel.QueueBind(
		q.Name, routingKey, exchange, false, nil,
	)
}

// Esto manda un mensaje (en JSON) al exchange con la clave "activities.created".
func (r *Rabbit) Publish(exchange, routingKey string, body []byte) error {
	return r.Channel.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

// Libera recursos cuando el programa termina o el microservicio se apaga.
func (r *Rabbit) Close() {
	if err := r.Channel.Close(); err != nil {
		log.Println("Error cerrando canal:", err)
	}
	if err := r.Conn.Close(); err != nil {
		log.Println("Error cerrando conexión:", err)
	}
}
