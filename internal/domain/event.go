package domain

import "time"

type Event struct {
	Id        int               `json:"id"`
	EventType string            `json:"event_type"`
	Data      map[string]string `json:"data"`
	Status    string            `json:"status"`
	CreatedAt time.Time         `json:"created_at"`
	SentAt    time.Time         `json:"sent_at"`
}

func NewPostEvent(
	EventType string,
	PostName string,
	PostId string,
	PostText string,
	AuthorId string,
) Event {
	return Event{
		EventType: EventType,
		Data: map[string]string{
			"Id":        PostId,
			"Name":      PostName,
			"Author_id": AuthorId,
			"Text":      PostText,
		},
	}
}
