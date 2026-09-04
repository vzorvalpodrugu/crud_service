package repository

import (
	"context"
	"crud_service/internal/domain"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &userRepository{pool: pool}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	query := `
        INSERT INTO users (name, email)
        VALUES ($1, $2)
        RETURNING id, name, email, created_at, updated_at
    `

	err := r.pool.QueryRow(ctx, query,
		user.Name,
		user.Email,
	).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.Created_at,
		&user.Updated_at,
	)
	if err != nil {
		return nil, fmt.Errorf("UserRepository.Create: %w", err)
	}

	return user, nil
}

func (r *userRepository) GetById(ctx context.Context, id int) (*domain.User, error) {
	query := `
		SELECT Id, Name, Email, Created_at, Updated_at 
		FROM users 
		WHERE Id = $1
	`

	user := &domain.User{}

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.Created_at,
		&user.Updated_at,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", ErrNotFound)
		}
		return nil, fmt.Errorf("UserRepository.GetById: %w", err)
	}

	return user, nil
}

func (r *userRepository) GetAll(ctx context.Context) ([]*domain.User, error) {
	query := `
	SELECT Id, Name, Email, Created_at, Updated_at
	FROM users
	ORDER BY Id DESC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("userRepository.GetAll: %w", err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		user := &domain.User{}
		if err := rows.Scan(
			&user.Id,
			&user.Name,
			&user.Email,
			&user.Created_at,
			&user.Updated_at,
		); err != nil {
			return nil, fmt.Errorf("userRepository.GetAll scan: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("userRepository.GetAll rows: %w", err)
	}

	return users, nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users
		SET Name = $1, Email = $2
		WHERE Id = $3
	`

	result, err := r.pool.Exec(ctx, query,
		user.Name,
		user.Email,
		user.Id,
	)

	if err != nil {
		return fmt.Errorf("userRepository.Update: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("User not found: %w", ErrNotFound)
	}

	return nil
}
func (r *userRepository) Delete(ctx context.Context, id int) error {
	query := `
		DELETE FROM users WHERE Id = $1
	`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("userRepository.Delete: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("User not found: %w", ErrNotFound)
	}

	return nil
}
