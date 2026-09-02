package domain

import "time"

type Comment struct {
	Id         int
	Author_id  int
	Post_id    int
	Text       string
	Created_at time.Time
	Updated_at time.Time
}
