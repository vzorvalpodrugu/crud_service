package kafka

import (
	"context"
	"fmt"

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
		Async:        true,
		//BatchTimeout: 10 * time.Millisecond,
		BatchSize: 100,
	}

	return &Producer{writer: writer}
}

type Message struct {
	Topic     string
	EventType string
	Key       string
	Value     []byte
}

func (p *Producer) Publish(ctx context.Context, msgs ...kafka.Message) error {
	//err := p.writer.WriteMessages(ctx, kafka.Message{
	//	Topic: topic,
	//	Key:   []byte(msg.Key),
	//	Value: msg.Value,
	//	Headers: []kafka.Header{
	//		{
	//			Key:   "event_type",
	//			Value: []byte(msg.EventType),
	//		},
	//	},
	//})

	err := p.writer.WriteMessages(ctx, msgs...)

	if err != nil {
		return fmt.Errorf("kafka.Producer.Publish: %w", err)
	}

	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
