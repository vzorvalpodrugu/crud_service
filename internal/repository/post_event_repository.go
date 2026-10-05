package repository

import (
	"context"
	"crud_service/internal/domain"
	"fmt"
	"log"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type postEventRepository struct {
	conn driver.Conn
}

func NewPostEventRepository(conn driver.Conn) PostEventRepository {
	return &postEventRepository{conn: conn}
}

func (p *postEventRepository) SaveOne(ctx context.Context, event *domain.PostEvent) (*domain.PostEvent, error) {
	query := `
		INSERT INTO analytics.posts_events(
			event_id,   
			event_type, 
			post_id, 
			post_name, 
			post_text, 
			post_author_id, 
			post_created_at, 
			post_updated_at,                                  
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);
	`

	if err := p.conn.Exec(ctx, query,
		event.EventId,
		event.EventType,
		event.PostId,
		event.PostName,
		event.PostText,
		event.PostAuthorId,
		event.PostCreatedAt,
		event.PostUpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("PostEventRepository.Save Exec: %w", err)
	}

	log.Printf("ClickHouse has successful saved event")

	return event, nil

}

func (p *postEventRepository) SaveBatch(ctx context.Context, events []*domain.PostEvent) error {
	query := `
		INSERT INTO analytics.posts_events(
			event_id,   
			event_type, 
			post_id, 
			post_name, 
			post_text, 
			post_author_id, 
			post_created_at, 
			post_updated_at,                                  
		)
	`
	batch, err := p.conn.PrepareBatch(ctx, query)
	if err != nil {
		return fmt.Errorf("PostEventRepository.SaveBatch PrepareBatch: %w", err)
	}

	for id, row := range events {
		if err := batch.Append(
			row.EventId,
			row.EventType,
			row.PostId,
			row.PostName,
			row.PostText,
			row.PostAuthorId,
			row.PostCreatedAt,
			row.PostUpdatedAt,
		); err != nil {
			return fmt.Errorf("PostEventRepository.SaveBatch row %d Append: %w", id, err)
		}
	}

	if err := batch.Send(); err != nil {
		return fmt.Errorf("PostEventRepository.SaveBatch Send: %w", err)
	}

	return nil
}
