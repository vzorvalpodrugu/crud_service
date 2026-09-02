package repository

import (
	"context"
	"crud_service/internal/domain"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type commentRepository struct {
	pool *pgxpool.Pool
}

func NewCommentRepository(pool *pgxpool.Pool) CommentRepository {
	return &commentRepository{pool: pool}
}

func (r *commentRepository) Create(ctx context.Context, comment *domain.Comment) (*domain.Comment, error) {
	query := `
		INSERT INTO comments (Author_id, Post_id, Text)
		VALUES ($1, $2, $3)
	`

	_, err := r.pool.Exec(ctx, query,
		comment.Author_id,
		comment.Post_id,
		comment.Text,
	)
	if err != nil {
		return nil, fmt.Errorf("CommentRepository.Create: %w", err)
	}

	return comment, nil
}

func (r *commentRepository) GetById(ctx context.Context, id int) (*domain.Comment, error) {
	query := `
			SELECT Id, Author_id, Post_id, Text, Created_at, Updated_at
			FROM comments
			WHERE Id = $1
		`

	comment := &domain.Comment{}

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&comment.Id,
		&comment.Author_id,
		&comment.Post_id,
		&comment.Text,
		&comment.Created_at,
		&comment.Updated_at,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("comment not found: %w", ErrNotFound)
		}

		return nil, fmt.Errorf("CommentRepository.GetById: %w", err)
	}

	return comment, nil
}

func (r *commentRepository) GetAll(ctx context.Context) ([]*domain.Comment, error) {
	query := `
		SELECT Id, Author_id, Post_id, Text, Created_at, Updated_at
		FROM comments
		ORDER BY Id DESC
	`

	var comments []*domain.Comment

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("CommentRepository.GetAll: %w", err)
	}

	for rows.Next() {
		comment := &domain.Comment{}

		err := rows.Scan(
			&comment.Id,
			&comment.Author_id,
			&comment.Post_id,
			&comment.Text,
			&comment.Created_at,
			&comment.Updated_at,
		)
		if err != nil {
			return nil, fmt.Errorf("CommentRepository.GetById scan: %w", err)
		}

		comments = append(comments, comment)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("CommentRepository.GetById rows: %w", err)
	}

	return comments, nil
}

func (r *commentRepository) Update(ctx context.Context, comment *domain.Comment) error {
	query := `
		UPDATE comments
		SET (Author_id = $1, Post_id = $2, Text = $3)
		WHERE Id = $4
	`

	result, err := r.pool.Exec(ctx, query,
		comment.Author_id,
		comment.Post_id,
		comment.Text,
		comment.Id,
	)
	if err != nil {
		return fmt.Errorf("CommentRepository.Update: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("comment not found: %w", ErrNotFound)
	}

	return nil
}

func (r *commentRepository) Delete(ctx context.Context, id int) error {
	query := `
		DELETE FROM comments WHERE Id = $1
	`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("CommentRepository.Delete: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("comment not found: %w", ErrNotFound)
	}

	return nil
}
