package domain

import "time"

type Post struct {
	id         int
	name       string
	author_id  int
	text       string
	created_at time.Time
	updated_at time.Time
}
