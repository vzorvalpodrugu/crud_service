package outbox

import (
	"context"
	"crud_service/internal/domain"
	"crud_service/internal/kafka"
	"crud_service/internal/repository"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type Processor struct {
	outboxRepo   repository.OutboxRepository
	producer     *kafka.Producer
	pollInterval time.Duration
}

func NewProcessor(
	outboxRepo repository.OutboxRepository,
	producer *kafka.Producer,
) *Processor {
	return &Processor{
		outboxRepo:   outboxRepo,
		producer:     producer,
		pollInterval: 5 * time.Second,
	}
}

// Горутина старта
func (p *Processor) Start(ctx context.Context) {
	go func() {
		log.Println("Outbox processor started")

		ticker := time.NewTicker(p.pollInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Println("Outbox processor stopped")
				return
			case <-ticker.C:
				if err := p.process(ctx); err != nil {
					log.Println("Outbox processor error: %v", err)
				}
			}
		}
	}()
}

// Получает необработанные события и вызывает для них handleEvent и MarkSent, остальные помечает failed
func (p *Processor) process(ctx context.Context) error {
	events, err := p.outboxRepo.GetPending(ctx)
	if err != nil {
		return fmt.Errorf("processor.process GetPending: %w", err)
	}

	if len(events) == 0 {
		return nil
	}

	log.Println("Outbox processor: fount %d pending events", len(events))

	for _, event := range events {
		if err := p.handleEvent(ctx, event); err != nil {
			log.Printf("Outbox processor: failed to handle event %d: %v", event.Id, err)

			if err := p.outboxRepo.MarkFailed(ctx, event.Id); err != nil {
				log.Printf("Outbox processor: failed to mark event %d as failed: %v", event.Id, err)
			}
			continue
		}

		if err := p.outboxRepo.MarkSent(ctx, event.Id); err != nil {
			log.Printf("Outbox processor: failed to mark event %d as sent: %v", event.Id, err)
		}

	}

	return nil
}

func (p *Processor) handleEvent(ctx context.Context, event *domain.OutboxEvent) error {
	topic := p.resolveTopic(event.EventType)
	if topic == "" {
		return fmt.Errorf("unknown event type: %s", event.EventType)
	}

	key, err := p.extractKey(event.Data)
	if err != nil {
		return fmt.Errorf("failed to extract key: %w", err)
	}

	return p.producer.Publish(ctx, topic, kafka.Message{
		EventType: event.EventType,
		Key:       key,
		Value:     event.Data,
	})
}

func (p *Processor) resolveTopic(eventType string) string {
	switch eventType {
	case domain.EventPostCreated,
		domain.EventPostUpdated,
		domain.EventPostDeleted:
		return kafka.TopicPostsEvents
	default:
		return ""
	}

}

func (p *Processor) extractKey(data []byte) (string, error) {
	var payload map[string]any

	//log.Printf(string(data))

	if err := json.Unmarshal(data, &payload); err != nil {
		return "", fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	id, ok := payload["Id"].(float64)

	//log.Printf("data_id: %f", id)

	if !ok {
		return "", fmt.Errorf("id not found in payload")
	}

	return fmt.Sprintf("%d", int(id)), nil
}
