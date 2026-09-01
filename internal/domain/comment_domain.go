package domain

import "time"

type Comment struct {
	id         int
	author_id  int
	post_id    int
	text       string
	created_at time.Time
	updated_at time.Time
}
