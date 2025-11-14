package services

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// PublishMessage representa un mensaje a publicar
type PublishMessage struct {
	Action string
	ID     string
}

// PublishQueue implementa una cola con workers que reintentan publishes
type PublishQueue struct {
	underlying ActivityPublisher
	queue      chan PublishMessage
	maxRetries int
	backoff    time.Duration

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// NewPublishQueue crea una nueva instancia
func NewPublishQueue(underlying ActivityPublisher, queueSize, maxRetries int, backoff time.Duration) *PublishQueue {
	return &PublishQueue{
		underlying: underlying,
		queue:      make(chan PublishMessage, queueSize),
		maxRetries: maxRetries,
		backoff:    backoff,
	}
}

// Start lanza N workers que procesan la cola. Debe llamarse antes de usar Publish.
func (q *PublishQueue) Start(parentCtx context.Context, workers int) {
	q.ctx, q.cancel = context.WithCancel(parentCtx)
	for i := 0; i < workers; i++ {
		q.wg.Add(1)
		go q.worker(i)
	}
}

// Stop cancela el contexto y espera a que los workers terminen.
func (q *PublishQueue) Stop() {
	if q.cancel != nil {
		q.cancel()
	}
	q.wg.Wait()
}

// Publish encola el mensaje para publicación. Respeta el context pasado.
func (q *PublishQueue) Publish(ctx context.Context, action string, activityID string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	msg := PublishMessage{Action: action, ID: activityID}

	select {
	case q.queue <- msg:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-q.ctx.Done():
		return fmt.Errorf("publish queue is stopping")
	}
}

func (q *PublishQueue) worker(idx int) {
	defer q.wg.Done()
	for {
		select {
		case <-q.ctx.Done():
			return
		case msg, ok := <-q.queue:
			if !ok {
				return
			}
			q.process(msg)
		}
	}
}

func (q *PublishQueue) process(msg PublishMessage) {
	var lastErr error
	for attempt := 1; attempt <= q.maxRetries; attempt++ {
		// usar un contexto corto por intento
		attemptCtx, cancel := context.WithTimeout(q.ctx, 5*time.Second)
		err := q.underlying.Publish(attemptCtx, msg.Action, msg.ID)
		cancel()
		if err == nil {
			fmt.Printf("✅ Published event: action=%s, activity_id=%s\n", msg.Action, msg.ID)
			return
		}
		lastErr = err
		// backoff exponencial simple
		time.Sleep(q.backoff * time.Duration(attempt))
	}
	// último intento fallido: log y continuar
	fmt.Printf("❌ Publish failed after retries for id=%s action=%s: %v\n", msg.ID, msg.Action, lastErr)
}
