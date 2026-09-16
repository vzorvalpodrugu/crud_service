package domain

import "time"

type OutboxEvent struct {
	Id        int        `json:"id"`
	EventType string     `json:"event_type"`
	Data      []byte     `json:"data"`
	Status    string     `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	SentAt    *time.Time `json:"sent_at"`
}

const (
	EventPostCreated = "post.created"
	EventPostUpdated = "post.updated"
	EventPostDeleted = "post.deleted"
)

const (
	OutboxStatusPending = "PENDING"
	OutboxStatusSent    = "SENT"
	OutboxStatusFailed  = "FAILED"
)
