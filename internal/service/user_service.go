package service

import (
	"context"
	"fmt"

	"crud_service/internal/domain"
	"crud_service/internal/repository"
)

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) Create(ctx context.Context, name, email string) (*domain.User, error) {
	user := &domain.User{
		Name:  name,
		Email: email,
	}

	created, err := s.userRepo.Create(ctx, user)

	if err != nil {
		return nil, fmt.Errorf("userService.Create: %w", err)
	}

	return created, nil
}

func (s *userService) GetById(ctx context.Context, id int) (*domain.User, error) {
	user, err := s.userRepo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("userService.GetByID: %w", err)
	}
	return user, nil
}

func (s *userService) GetAll(ctx context.Context) ([]*domain.User, error) {
	users, err := s.userRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("userService.GetAll: %w", err)
	}
	return users, nil
}

func (s *userService) Update(ctx context.Context, name, email string, id int) error {
	user, err := s.userRepo.GetById(ctx, id)
	if err != nil {
		return fmt.Errorf("userService.Update: %w", err)
	}

	if name != "" {
		user.Name = name
	}
	if email != "" {
		user.Email = email
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("userService.Update: %w", err)
	}

	return nil
}

func (s *userService) Delete(ctx context.Context, id int) error {
	if err := s.userRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("userService.Delete: %w", err)
	}
	return nil
}
