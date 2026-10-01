package outbox

import (
	"context"
	"crud_service/internal/domain"
	"crud_service/internal/kafka"
	"crud_service/internal/metrics"
	"crud_service/internal/repository"
	"encoding/json"
	"fmt"
	"log"
	"time"

	kafka2 "github.com/segmentio/kafka-go"
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
		pollInterval: 1 * time.Millisecond,
	}
}

// Горутина старта
func (p *Processor) Start(ctx context.Context) {
	go func() {
		log.Println("Outbox processor started")

		counter := 0
		ticker := time.NewTicker(p.pollInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Println("Outbox processor stopped")
				return
			case <-ticker.C:
				go func() {
					//start := time.Now()
					if err := p.process(ctx, counter); err != nil {
						log.Println("Outbox processor error: %v", err)
					}
					//dur := time.Since(start)
					//log.Printf("ВРЕМЯ ВЫПОЛНЕНИЯ PROCESS ПОД НОМЕРОМ %d ЗАНЯЛО %v", counter, dur)
				}()
				counter++
			}
		}
	}()
}

// Получает необработанные события и вызывает для них handleEvent и MarkSent, остальные помечает failed
func (p *Processor) process(ctx context.Context, counter int) error {
	events, err := p.outboxRepo.GetPending(ctx)
	if err != nil {
		return fmt.Errorf("processor.process GetPending: %w", err)
	}
	//log.Printf("PROCESS ПОД НОМЕРОМ %d ВЗЯЛ НА СЕБЯ %d ЗАДАЧ", counter, int(len(events)))
	if len(events) == 0 {
		return nil
	}

	//log.Println("Outbox %d processor: fount %d pending events", counter, len(events))

	count, _ := p.outboxRepo.GetCountPendingEvents(ctx)
	metrics.OutboxPendingEvents.Set(float64(count))

	var kafkaMessages []kafka2.Message

	for _, event := range events {
		kafkaEvent, err := p.handleEvent(event)
		if err != nil {
			log.Printf("Outbox processor: failed to handle event %d: %v", event.Id, err)

			if err := p.outboxRepo.MarkFailed(ctx, event.Id); err != nil {
				log.Printf("Outbox processor: failed to mark event %d as failed: %v", event.Id, err)
			}
			continue
		}

		kafkaMessages = append(kafkaMessages, kafkaEvent)
		if err := p.outboxRepo.MarkSent(ctx, event.Id); err != nil {
			log.Printf("Outbox processor: failed to mark event %d as sent: %v", event.Id, err)
		}

	}

	//start := time.Now()
	if err := p.producer.Publish(ctx, kafkaMessages...); err != nil {
		log.Printf("Outbox processor: failed to publish %d events : %v", count, err)
	}
	//end := time.Since(start)

	//log.Printf("ВРЕМЯ ВЫПОЛНЕНИЯ PUBLISH НОМЕР %d ЗАНЯЛО %v", counter, end)

	count, _ = p.outboxRepo.GetCountPendingEvents(ctx)
	metrics.OutboxPendingEvents.Set(float64(count))

	return nil
}

func (p *Processor) handleEvent(event *domain.OutboxEvent) (kafka2.Message, error) {
	topic := p.resolveTopic(event.EventType)
	if topic == "" {
		return kafka2.Message{}, fmt.Errorf("unknown event type: %s", event.EventType)
	}

	key, err := p.extractKey(event.Data)
	if err != nil {
		return kafka2.Message{}, fmt.Errorf("failed to extract key: %w", err)
	}

	kafkaEvent := kafka2.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: event.Data,
		Headers: []kafka2.Header{
			{
				Key:   "event_type",
				Value: []byte(event.EventType),
			},
		},
	}
	metrics.KafkaProducedTotal.WithLabelValues(event.EventType).Inc()
	return kafkaEvent, nil

	//if err != nil {
	//	metrics.KafkaProducedErrorsTotal.WithLabelValues(event.EventType).Inc()
	//}
	//
	//metrics.KafkaProducedTotal.WithLabelValues(event.EventType).Inc()
	//
	//return nil
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
