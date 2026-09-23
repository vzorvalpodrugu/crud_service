package kafka

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(brokers []string) Consumer {
	return Consumer{reader: kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: "posts.group",
		Topic:   TopicPostsEvents,
	})}
}

func (r *Consumer) Read(ctx context.Context) (kafka.Message, error) {
	//log.Printf("Try to read message")
	msg, err := r.reader.ReadMessage(ctx)
	//log.Printf("ok i get a message")
	if err != nil {
		return kafka.Message{}, fmt.Errorf("Consumer.Read ReadMessage: %w", err)
	}

	return msg, nil
}

func (r *Consumer) Close() error {
	if err := r.reader.Close(); err != nil {
		return fmt.Errorf("Consumer.Close close: %w", err)
	}

	return nil
}
