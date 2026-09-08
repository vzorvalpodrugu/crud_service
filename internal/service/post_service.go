package service

import (
	"context"
	"crud_service/internal/cache"
	"fmt"
	"log"

	"crud_service/internal/domain"
	"crud_service/internal/repository"
)

type postService struct {
	postRepo  repository.PostRepository
	postCache cache.PostCache
}

func NewPostService(postRepo repository.PostRepository, postCache cache.PostCache) PostService {
	return &postService{
		postRepo:  postRepo,
		postCache: postCache,
	}
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

	if err = s.postCache.SetById(ctx, created); err != nil {
		log.Println("PostService.Create SetById cache failed")
	}
	return created, nil
}

func (s *postService) GetById(ctx context.Context, id int) (*domain.Post, error) {
	post, err := s.postCache.GetById(ctx, id)
	if err != nil {
		log.Println("PostService.GetById cache failed")
	}

	if post != nil {
		return post, nil
	}

	post, err = s.postRepo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("postService.GetByID: %w", err)
	}

	if err = s.postCache.SetById(ctx, post); err != nil {
		log.Println("PostService.GetById SetById cache failed")
	}

	return post, nil
}

func (s *postService) GetAll(ctx context.Context) ([]*domain.Post, error) {
	posts, err := s.postCache.GetAll(ctx)
	if err != nil {
		log.Println("PostService.GetAll cache failed")
	}

	if posts != nil {
		return posts, nil
	}

	posts, err = s.postRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("postService.GetAll: %w", err)
	}

	if err = s.postCache.SetAll(ctx, posts); err != nil {
		log.Println("PostService.GetAll SetAll cache failed")
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

	if err = s.postCache.Invalidate(ctx, post.Id); err != nil {
		log.Println("PostService.Update Invalidate cache failed")
	}

	if err = s.postCache.SetById(ctx, post); err != nil {
		log.Println("PostService.Update SetById cache failed")
	}

	return nil
}

func (s *postService) Delete(ctx context.Context, id int) error {
	if err := s.postRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("postService.Delete: %w", err)
	}

	if err := s.postCache.Invalidate(ctx, id); err != nil {
		log.Println("PostService.Delete Invalidate cache failed")
	}

	return nil
}
