package repository

import (
	"context"
	"crud_service/internal/domain"

	"github.com/jackc/pgx/v5"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	GetById(ctx context.Context, id int) (*domain.User, error)
	GetAll(ctx context.Context) ([]*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id int) error
	AssignRole(ctx context.Context, userId int, role string) error
	RemoveRole(ctx context.Context, userId int, role string) error
	GetRoles(ctx context.Context, userId int) ([]string, error)
}

type PostRepository interface {
	Create(ctx context.Context, tx pgx.Tx, post *domain.Post) (*domain.Post, error)
	GetById(ctx context.Context, id int) (*domain.Post, error)
	GetAll(ctx context.Context) ([]*domain.Post, error)
	Update(ctx context.Context, post *domain.Post) error
	Delete(ctx context.Context, id int) error
}

type CommentRepository interface {
	Create(ctx context.Context, comment *domain.Comment) (*domain.Comment, error)
	GetById(ctx context.Context, id int) (*domain.Comment, error)
	GetAll(ctx context.Context) ([]*domain.Comment, error)
	Update(ctx context.Context, comment *domain.Comment) error
	Delete(ctx context.Context, id int) error
}

type OutboxRepository interface {
	Create(ctx context.Context, tx pgx.Tx, event *domain.OutboxEvent) error
	GetPending(ctx context.Context) ([]*domain.OutboxEvent, error)
	MarkSent(ctx context.Context, id int) error
	MarkFailed(ctx context.Context, id int) error
}
