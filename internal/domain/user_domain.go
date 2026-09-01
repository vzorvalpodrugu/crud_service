package domain

import "time"

type User struct {
	Id         int
	Name       string
	Email      string
	Created_at time.Time
	Updated_at time.Time
}
