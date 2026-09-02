package domain

import "time"

type Post struct {
	Id         int
	Name       string
	Author_id  int
	Text       string
	Created_at time.Time
	Updated_at time.Time
}
