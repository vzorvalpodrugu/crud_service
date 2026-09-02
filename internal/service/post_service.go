package service

import (
	"context"
	"fmt"

	"crud_service/internal/domain"
	"crud_service/internal/repository"
)

type postService struct {
	postRepo repository.PostRepository
}

func NewPostService(postRepo repository.PostRepository) PostService {
	return &postService{postRepo: postRepo}
}

func (s *postService) Create(ctx context.Context, name, text string, author_id int) (*domain.Post, error) {
	post := &domain.Post{
		Name:      name,
		Text:      text,
		Author_id: author_id,
	}

	created, err := s.postRepo.Create(ctx, post)

	if err != nil {
		return nil, fmt.Errorf("postService.Create: %w", err)
	}

	return created, nil
}

func (s *postService) GetById(ctx context.Context, id int) (*domain.Post, error) {
	post, err := s.postRepo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("postService.GetByID: %w", err)
	}
	return post, nil
}

func (s *postService) GetAll(ctx context.Context) ([]*domain.Post, error) {
	posts, err := s.postRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("postService.GetAll: %w", err)
	}
	return posts, nil
}

func (s *postService) Update(ctx context.Context, name, text string, id, author_id int) error {
	post, err := s.postRepo.GetById(ctx, id)
	if err != nil {
		return fmt.Errorf("postService.Update: %w", err)
	}

	if name != "" {
		post.Name = name
	}
	if text != "" {
		post.Text = text
	}
	if author_id >= 0 {
		post.Author_id = author_id
	}

	if err := s.postRepo.Update(ctx, post); err != nil {
		return fmt.Errorf("postService.Update: %w", err)
	}

	return nil
}

func (s *postService) Delete(ctx context.Context, id int) error {
	if err := s.postRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("postService.Delete: %w", err)
	}
	return nil
}
