package kafka

import (
	"context"
	"fmt"
	"log"
	"time"
)

type Processor struct {
	consumer     *Consumer
	pollInterval time.Duration
}

func NewProcessor(consumer *Consumer) Processor {
	return Processor{
		consumer:     consumer,
		pollInterval: 10 * time.Millisecond,
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

	log.Printf("Message is %s", string(msg.Value))

	return nil
}
