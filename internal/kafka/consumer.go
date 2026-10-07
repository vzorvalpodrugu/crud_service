package kafka

import (
	"github.com/segmentio/kafka-go"
)

type KafkaConsumerGroup struct {
	readers []*kafka.Reader
}

func NewConsumerGroup(brokers []string, num int) KafkaConsumerGroup {
	var readers []*kafka.Reader

	for _ = range num {
		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			GroupID: "posts.group",
			Topic:   TopicPostsEvents,
		})
		readers = append(readers, reader)
	}
	return KafkaConsumerGroup{readers: readers}
}

func (r *KafkaConsumerGroup) Close() {
	for idx := range len(r.readers) {
		r.readers[idx].Close()
	}
}
