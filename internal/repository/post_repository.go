package repository

import (
	"context"
	"crud_service/internal/domain"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postRepository struct {
	pool *pgxpool.Pool
}

func NewPostRepository(pool *pgxpool.Pool) PostRepository {
	return &postRepository{pool: pool}
}

func (r *postRepository) Create(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	query := `
		INSERT INTO posts (Name, Author_id, Text)
		VALUES ($1, $2, $3)
	`

	_, err := r.pool.Exec(ctx, query,
		post.Name,
		post.Author_id,
		post.Text,
	)
	if err != nil {
		return nil, fmt.Errorf("PostRepository.Create: %w", err)
	}

	return post, nil
}

func (r *postRepository) GetById(ctx context.Context, id int) (*domain.Post, error) {
	query := `
			SELECT Id, Name, Author_id, Text, Created_at, Updated_at
			FROM posts
			WHERE Id = $1
		`

	post := &domain.Post{}

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&post.Id,
		&post.Name,
		&post.Author_id,
		&post.Text,
		&post.Created_at,
		&post.Updated_at,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("post not found: %w", ErrNotFound)
		}

		return nil, fmt.Errorf("PostRepository.GetById: %w", err)
	}

	return post, nil
}

func (r *postRepository) GetAll(ctx context.Context) ([]*domain.Post, error) {
	query := `
		SELECT Id, Name, Author_id, Text, Created_at, Updated_at
		FROM posts
		ORDER BY Id DESC
	`

	var posts []*domain.Post

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("PostRepository.GetAll: %w", err)
	}

	for rows.Next() {
		post := &domain.Post{}

		err := rows.Scan(
			&post.Id,
			&post.Name,
			&post.Author_id,
			&post.Text,
			&post.Created_at,
			&post.Updated_at,
		)
		if err != nil {
			return nil, fmt.Errorf("PostRepository.GetById scan: %w", err)
		}

		posts = append(posts, post)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("PostRepository.GetById rows: %w", err)
	}

	return posts, nil
}

func (r *postRepository) Update(ctx context.Context, post *domain.Post) error {
	query := `
		UPDATE posts
		SET (Name = $1, Author_id = $2, Text = $3)
		WHERE Id = $4
	`

	result, err := r.pool.Exec(ctx, query,
		post.Name,
		post.Author_id,
		post.Text,
		post.Id,
	)
	if err != nil {
		return fmt.Errorf("PostRepository.Update: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("post not found: %w", ErrNotFound)
	}

	return nil
}

func (r *postRepository) Delete(ctx context.Context, id int) error {
	query := `
		DELETE FROM posts WHERE Id = $1
	`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("PostRepository.Delete: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("post not found: %w", ErrNotFound)
	}

	return nil
}
