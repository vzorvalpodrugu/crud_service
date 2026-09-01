package repository

import (
	"context"
	"crud_service/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	GetById(ctx context.Context, id int) (*domain.User, error)
	GetAll(ctx context.Context) ([]*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id int) error
}

type PostRepository interface {
	Create(ctx context.Context, post *domain.Post) (*domain.Post, error)
	GetById(ctx context.Context, id int) (*domain.Post, error)
	GetAll(ctx context.Context) ([]*domain.Post, error)
	Update(ctx context.Context, post *domain.Post) (*domain.Post, error)
	Delete(ctx context.Context, id int) error
}

type CommentRepository interface {
	Create(ctx context.Context, comment *domain.Comment) (*domain.Comment, error)
	GetById(ctx context.Context, id int) (*domain.Comment, error)
	GetAll(ctx context.Context) ([]*domain.Comment, error)
	Update(ctx context.Context, comment *domain.Comment) (*domain.Comment, error)
	Delete(ctx context.Context, id int) error
}
