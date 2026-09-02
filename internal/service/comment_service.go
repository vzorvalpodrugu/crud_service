package service

import (
	"context"
	"fmt"

	"crud_service/internal/domain"
	"crud_service/internal/repository"
)

type commentService struct {
	commentRepo repository.CommentRepository
}

func NewCommentService(commentRepo repository.CommentRepository) CommentService {
	return &commentService{commentRepo: commentRepo}
}

func (s *commentService) Create(ctx context.Context, author_id, post_id int, text string) (*domain.Comment, error) {
	comment := &domain.Comment{
		Author_id: author_id,
		Post_id:   post_id,
		Text:      text,
	}

	created, err := s.commentRepo.Create(ctx, comment)

	if err != nil {
		return nil, fmt.Errorf("commentService.Create: %w", err)
	}

	return created, nil
}

func (s *commentService) GetById(ctx context.Context, id int) (*domain.Comment, error) {
	comment, err := s.commentRepo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("userService.GetByID: %w", err)
	}
	return comment, nil
}

func (s *commentService) GetAll(ctx context.Context) ([]*domain.Comment, error) {
	comments, err := s.commentRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("commentService.GetAll: %w", err)
	}
	return comments, nil
}

func (s *commentService) Update(ctx context.Context, id, author_id, post_id int, text string) error {
	comment, err := s.commentRepo.GetById(ctx, id)
	if err != nil {
		return fmt.Errorf("commentService.Update: %w", err)
	}

	if author_id >= 0 {
		comment.Author_id = author_id
	}
	if post_id >= 0 {
		comment.Post_id = post_id
	}
	if text != "" {
		comment.Text = text
	}

	if err := s.commentRepo.Update(ctx, comment); err != nil {
		return fmt.Errorf("commentService.Update: %w", err)
	}

	return nil
}

func (s *commentService) Delete(ctx context.Context, id int) error {
	if err := s.commentRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("commentService.Delete: %w", err)
	}
	return nil
}
