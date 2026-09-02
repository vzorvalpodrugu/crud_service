package service

import (
	"context"
	"crud_service/internal/domain"
)

type UserService interface {
	Create(ctx context.Context, name string, email string) (*domain.User, error)
	GetById(ctx context.Context, id int) (*domain.User, error)
	GetAll(ctx context.Context) ([]*domain.User, error)
	Update(ctx context.Context, name, email string, id int) error
	Delete(ctx context.Context, id int) error
}

type PostService interface {
	Create(ctx context.Context, name, text string, author_id int) (*domain.Post, error)
	GetById(ctx context.Context, id int) (*domain.Post, error)
	GetAll(ctx context.Context) ([]*domain.Post, error)
	Update(ctx context.Context, name, text string, id, author_id int) error
	Delete(ctx context.Context, id int) error
}

type CommentService interface {
	Create(ctx context.Context, author_id, post_id int, text string) (*domain.Comment, error)
	GetById(ctx context.Context, id int) (*domain.Comment, error)
	GetAll(ctx context.Context) ([]*domain.Comment, error)
	Update(ctx context.Context, id, author_id, post_id int, text string) error
	Delete(ctx context.Context, id int) error
}
