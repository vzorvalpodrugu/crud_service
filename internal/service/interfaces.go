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
