package messaging

import (
	"encoding/json"

	"github.com/streadway/amqp"
)

// Publisher interface
type Publisher interface {
	Publish(operation string, entityID string) error
	Close() error
}

type rabbitPublisher struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	exchange string
}

// NewRabbitPublisher dials and declares exchange
func NewRabbitPublisher(url, exchange string) (Publisher, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}
	// fanout exchange
	if err := ch.ExchangeDeclare(exchange, "fanout", true, false, false, false, nil); err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}
	return &rabbitPublisher{conn: conn, channel: ch, exchange: exchange}, nil
}

func (r *rabbitPublisher) Publish(operation string, entityID string) error {
	body := map[string]string{"op": operation, "id": entityID}
	b, _ := json.Marshal(body)
	return r.channel.Publish(r.exchange, "", false, false, amqp.Publishing{ContentType: "application/json", Body: b})
}

func (r *rabbitPublisher) Close() error {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}
