package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	TopicPostsEvents = "posts.events"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string) *Producer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireOne,
		Async:        false,
		BatchTimeout: 10 * time.Millisecond,
	}

	return &Producer{writer: writer}
}

type Message struct {
	EventType string
	Key       string
	Value     []byte
}

func (p *Producer) Publish(ctx context.Context, topic string, msg Message) error {
	err := p.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Key:   []byte(msg.Key),
		Value: msg.Value,
		Headers: []kafka.Header{
			{
				Key:   "event_type",
				Value: []byte(msg.EventType),
			},
		},
	})

	if err != nil {
		return fmt.Errorf("kafka.Producer.Publish: %w", err)
	}

	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
