package domain

import "time"

type Role struct {
	Id   int
	Name string
}
type User struct {
	Id         int
	Name       string
	Email      string
	Roles      []Role
	Created_at time.Time
	Updated_at time.Time
}
