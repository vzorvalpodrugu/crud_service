package handler

import (
	"crud_service/api/api"
	"crud_service/internal/domain"
)

// domain -> HTTP Response
func userToResponse(u *domain.User) api.UserResponse {
	return api.UserResponse{
		CreatedAt: &u.Created_at,
		Email:     &u.Email,
		Id:        &u.Id,
		Name:      &u.Name,
		UpdatedAt: &u.Updated_at,
	}
}

func postToResponse(p *domain.Post) api.PostResponse {
	return api.PostResponse{
		AuthorId:  &p.Author_id,
		CreatedAt: &p.Created_at,
		Id:        &p.Id,
		Name:      &p.Name,
		Text:      &p.Text,
		UpdatedAt: &p.Updated_at,
	}
}

func commentToResponse(p *domain.Comment) api.CommentResponse {
	return api.CommentResponse{
		AuthorId:  &p.Author_id,
		CreatedAt: &p.Created_at,
		Id:        &p.Id,
		PostId:    &p.Post_id,
		Text:      &p.Text,
		UpdatedAt: &p.Updated_at,
	}
}
