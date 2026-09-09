package repository

import "github.com/jackc/pgx/v5/pgxpool"

type eventRepository struct {
	pool *pgxpool.Pool
}

//func NewEventRep
