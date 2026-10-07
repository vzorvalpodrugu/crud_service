package kafka

import (
	"context"
	"crud_service/internal/domain"
	"crud_service/internal/metrics"
	"crud_service/internal/repository"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type PostEventBatch struct {
	events []*domain.PostEvent
	length int
}
type PostPayload struct {
	Id         int       `json:"Id"`
	Name       string    `json:"Name"`
	Text       string    `json:"Text"`
	Author_id  int       `json:"Author_Id"`
	Created_at time.Time `json:"Created_At"`
	Updated_at time.Time `json:"Updated_At"`
}

type Processor struct {
	batch         *PostEventBatch
	consumerGroup *KafkaConsumerGroup
	postEventRepo repository.PostEventRepository
	pollInterval  time.Duration
}

func NewProcessor(
	consumerGroup *KafkaConsumerGroup,
	postEventRepo repository.PostEventRepository,
) Processor {
	return Processor{
		batch:         &PostEventBatch{[]*domain.PostEvent{}, 0},
		consumerGroup: consumerGroup,
		postEventRepo: postEventRepo,
		pollInterval:  1 * time.Millisecond,
	}
}

func (p *Processor) Start(ctx context.Context) {
	for readerId := range len(p.consumerGroup.readers) {
		go func(readerIdx int) {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					if err := p.process(ctx, readerIdx); err != nil {
						log.Printf("Consumer-process.Start process: %v", err)
						continue
					}

					p.batch.length++

					//log.Printf("consumer id:%d обработал сообщение; lenght = %d", readerIdx, p.batch.length)
				}
			}
		}(readerId)
	}
	go func() {
		for {
			if p.batch.length >= 1000 {
				start := time.Now()
				if err := p.postEventRepo.SaveBatch(ctx, p.batch.events); err != nil {
					return
				}

				p.batch.length = 0
				p.batch.events = []*domain.PostEvent{}

				dur := time.Since(start)
				log.Printf("\n\n\nBATCH SUCCESSFUL\nTAKE: %v\n\n\n", dur)
			}
		}
	}()
}

func (p *Processor) process(ctx context.Context, readerId int) error {
	//log.Printf("Try to get a message")
	//start1 := time.Now()
	msg, err := p.consumerGroup.readers[readerId].ReadMessage(ctx)
	//dur1 := time.Since(start1)

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

	metrics.KafkaConsumedTotal.WithLabelValues(eventType).Inc()

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
	p.batch.events = append(p.batch.events, event)

	//start2 := time.Now()
	//if _, err := p.postEventRepo.SaveOne(ctx, event); err != nil {
	//	metrics.ClickHouseWriteErrorsTotal.WithLabelValues(eventType).Inc()
	//	return fmt.Errorf("Consumer-Processor.process Save: %w", err)
	//}
	//dur2 := time.Since(start2)

	//log.Printf("КОНСЬЮМЕР %d\n--------------------\nЧИТАЛ %v\nЗАПИСЫВАЛ В КЛИКХАУС %v\n---------------------\n", readerId, dur1, dur2)

	return nil
}
