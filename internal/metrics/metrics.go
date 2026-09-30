package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	OutboxPendingEvents = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "outbox_pending_events",
		Help: "Current number of pending events in outbox table",
	})

	OutboxCreatedEvents = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "outbox_created_events",
		Help: "Current number of created events in outbox table",
	})

	KafkaProducedTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kafka_produced_messages_total",
			Help: "Total number of messages produced to Kafka",
		},
		[]string{"event_type"},
	)

	KafkaProducedErrorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kafka_produced_errors_total",
			Help: "Total number of errors when producing to Kafka",
		},
		[]string{"event_type"},
	)

	KafkaConsumedTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kafka_consumed_messages_total",
			Help: "Total number of messages consumed from Kafka",
		},
		[]string{"event_type"},
	)

	ClickHouseWriteErrorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "clickhouse_write_errors_total",
			Help: "Total number of errors when writing to ClickHouse",
		},
		[]string{"event_type"},
	)
)
