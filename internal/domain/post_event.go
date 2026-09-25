package domain

import "time"

type PostEvent struct {
	EventId   int
	EventType string

	PostId        int
	PostName      string
	PostText      string
	PostAuthorId  int
	PostCreatedAt time.Time
	PostUpdatedAt time.Time
	ReceivedAt    time.Time
}
