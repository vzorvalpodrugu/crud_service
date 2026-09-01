package domain

import "time"

type User struct {
	id         int
	name       string
	email      string
	created_at time.Time
	updated_at time.Time
}
