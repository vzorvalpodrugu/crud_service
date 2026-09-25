package kafka

import (
	"context"
	"crud_service/internal/domain"
	"crud_service/internal/repository"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type PostPayload struct {
	Id         int       `json:"Id"`
	Name       string    `json:"Name"`
	Text       string    `json:"Text"`
	Author_id  int       `json:"Author_Id"`
	Created_at time.Time `json:"Created_At"`
	Updated_at time.Time `json:"Updated_At"`
}

type Processor struct {
	consumer      *Consumer
	postEventRepo repository.PostEventRepository
	pollInterval  time.Duration
}

func NewProcessor(
	consumer *Consumer,
	postEventRepo repository.PostEventRepository,
) Processor {
	return Processor{
		consumer:      consumer,
		postEventRepo: postEventRepo,
		pollInterval:  10 * time.Millisecond,
	}
}

func (p Processor) Start(ctx context.Context) {
	log.Printf("Consumer processor start")
	go func() {
		ticker := time.NewTicker(p.pollInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Printf("Consumer process stopped")
				return
			case <-ticker.C:
				if err := p.process(ctx); err != nil {
					log.Printf("Consumer process failed: %s", err)
				}
			}
		}

	}()
}

func (p Processor) process(ctx context.Context) error {
	//log.Printf("Try to get a message")
	msg, err := p.consumer.Read(ctx)

	if err != nil {
		return fmt.Errorf("Consumer-Processor.process Read: %w", err)
	}

	if len(msg.Value) == 0 {
		log.Printf("Message is empty")
	}

	var eventType string
	for _, header := range msg.Headers {
		if header.Key == "event_type" {
			eventType = string(header.Value)
		}
	}

	var postPayload PostPayload
	if err := json.Unmarshal(msg.Value, &postPayload); err != nil {
		return fmt.Errorf("Consumer-Processor.process Unmarshal: %w", err)
	}

	event := &domain.PostEvent{
		EventType:     eventType,
		PostId:        postPayload.Id,
		PostName:      postPayload.Name,
		PostText:      postPayload.Text,
		PostAuthorId:  postPayload.Author_id,
		PostCreatedAt: postPayload.Created_at,
		PostUpdatedAt: postPayload.Updated_at,
		ReceivedAt:    time.Now().UTC(),
	}

	if _, err := p.postEventRepo.Save(ctx, event); err != nil {
		return fmt.Errorf("Consumer-Processor.process Save: %w", err)
	}

	log.Printf("processor: successful save")

	return nil
}
