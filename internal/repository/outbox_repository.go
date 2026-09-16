package repository

import (
	"context"
	"crud_service/internal/domain"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type outboxRepository struct {
	pool *pgxpool.Pool
}

func NewOutboxRepository(pool *pgxpool.Pool) OutboxRepository {
	return &outboxRepository{pool: pool}
}

func (r *outboxRepository) Create(ctx context.Context, tx pgx.Tx, event *domain.OutboxEvent) error {
	query := `
        INSERT INTO outbox (event_type, data, status)
        VALUES ($1, $2, $3)
    `

	_, err := tx.Exec(ctx, query,
		event.EventType,
		event.Data,
		domain.OutboxStatusPending,
	)
	if err != nil {
		return fmt.Errorf("outboxRepository.Create: %w", err)
	}
	return nil
}

func (r *outboxRepository) GetPending(ctx context.Context) ([]*domain.OutboxEvent, error) {
	query := `
        SELECT id, event_type, data, status, created_at, sent_at
        FROM outbox
        WHERE status = $1
        ORDER BY created_at ASC
        LIMIT 100
        FOR UPDATE SKIP LOCKED
    `

	rows, err := r.pool.Query(ctx, query, domain.OutboxStatusPending)
	if err != nil {
		return nil, fmt.Errorf("outboxRepository.GetPending: %w", err)
	}
	defer rows.Close()

	var events []*domain.OutboxEvent
	for rows.Next() {
		event := &domain.OutboxEvent{}
		if err := rows.Scan(
			&event.Id,
			&event.EventType,
			&event.Data,
			&event.Status,
			&event.CreatedAt,
			&event.SentAt,
		); err != nil {
			return nil, fmt.Errorf("outboxRepository.GetPending scan: %w", err)
		}
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("outboxRepository.GetPending rows: %w", err)
	}

	return events, nil
}

func (r *outboxRepository) MarkSent(ctx context.Context, id int) error {
	query := `
        UPDATE outbox
        SET status = $1, sent_at = $2
        WHERE id = $3
    `
	now := time.Now().UTC()
	_, err := r.pool.Exec(ctx, query, domain.OutboxStatusSent, now, id)
	if err != nil {
		return fmt.Errorf("outboxRepository.MarkSent: %w", err)
	}
	return nil
}

func (r *outboxRepository) MarkFailed(ctx context.Context, id int) error {
	query := `
        UPDATE outbox
        SET status = $1
        WHERE id = $2
    `
	_, err := r.pool.Exec(ctx, query, domain.OutboxStatusFailed, id)
	if err != nil {
		return fmt.Errorf("outboxRepository.MarkFailed: %w", err)
	}
	return nil
}
