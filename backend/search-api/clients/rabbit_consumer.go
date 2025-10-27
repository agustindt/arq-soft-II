package clients

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"search-api/models"
	"search-api/repository"

	"github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn   *amqp091.Connection
	ch     *amqp091.Channel
	repo   *repository.SolrRepository
	actCli *ActivitiesClient
	queue  string
}

func NewConsumer(rabbitURL, queue string, repo *repository.SolrRepository, actCli *ActivitiesClient) (*Consumer, error) {
	conn, err := amqp091.Dial(rabbitURL)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	_, err = ch.QueueDeclare(queue, false, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	return &Consumer{conn: conn, ch: ch, repo: repo, actCli: actCli, queue: queue}, nil
}

type event struct {
	Action     string `json:"action"`
	ActivityID string `json:"activity_id"`
	ReservaID  string `json:"reserva_id"` // compat con reservations-api si la usás
}

func (c *Consumer) Start(ctx context.Context) error {
	deliveries, err := c.ch.Consume(c.queue, "", true, false, false, false, nil)
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case d := <-deliveries:
				var ev event
				if err := json.Unmarshal(d.Body, &ev); err != nil {
					log.Printf("[consumer] bad json: %v", err)
					continue
				}
				id := ev.ActivityID
				if id == "" {
					id = ev.ReservaID
				}
				if id == "" {
					log.Printf("[consumer] missing id in event: %s", string(d.Body))
					continue
				}
				action := strings.ToLower(ev.Action)
				switch action {
				case "create", "update":
					act, err := c.actCli.GetByID(id)
					if err != nil {
						log.Printf("[consumer] fetch %s: %v", id, err)
						continue
					}
					if err := c.repo.Upsert(models.Activity(*act)); err != nil {
						log.Printf("[consumer] upsert: %v", err)
					}
				case "delete":
					if err := c.repo.DeleteByID(id); err != nil {
						log.Printf("[consumer] delete: %v", err)
					}
				default:
					log.Printf("[consumer] unknown action %q", ev.Action)
				}
			}
		}
	}()
	return nil
}

func (c *Consumer) Close() { _ = c.ch.Close(); _ = c.conn.Close() }
