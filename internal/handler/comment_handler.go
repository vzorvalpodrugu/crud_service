package handler

import (
	"context"
	"errors"

	"crud_service/api/api"
	"crud_service/internal/repository"
	"crud_service/internal/service"
)

type CommentHandler struct {
	commentService service.CommentService
}

func NewCommentHandler(commentService service.CommentService) *CommentHandler {
	return &CommentHandler{commentService: commentService}
}

func (h *CommentHandler) CreateComment(ctx context.Context, request api.CreateCommentRequestObject) (api.CreateCommentResponseObject, error) {
	comment, err := h.commentService.Create(ctx, request.Body.AuthorId, request.Body.PostId, request.Body.Text)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return api.CreateComment404JSONResponse{
				NotFoundJSONResponse: api.NotFoundJSONResponse{
					Error: stringPtr("user or post not found"),
				},
			}, nil
		}
		return api.CreateComment500JSONResponse{
			InternalErrorJSONResponse: api.InternalErrorJSONResponse{
				Error: stringPtr(err.Error()),
			},
		}, nil
	}
	return api.CreateComment201JSONResponse(commentToResponse(comment)), nil
}

func (h *CommentHandler) ListComments(ctx context.Context, request api.ListCommentsRequestObject) (api.ListCommentsResponseObject, error) {
	comments, err := h.commentService.GetAll(ctx)
	if err != nil {
		return api.ListComments500JSONResponse{
			InternalErrorJSONResponse: api.InternalErrorJSONResponse{
				Error: stringPtr(err.Error()),
			},
		}, nil
	}

	resp := make([]api.CommentResponse, 0, len(comments))
	for _, c := range comments {
		resp = append(resp, commentToResponse(c))
	}
	return api.ListComments200JSONResponse(resp), nil
}

func (h *CommentHandler) GetCommentById(ctx context.Context, request api.GetCommentByIdRequestObject) (api.GetCommentByIdResponseObject, error) {
	comment, err := h.commentService.GetById(ctx, request.Id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return api.GetCommentById404JSONResponse{
				NotFoundJSONResponse: api.NotFoundJSONResponse{
					Error: stringPtr("comment not found"),
				},
			}, nil
		}
		return api.GetCommentById500JSONResponse{
			InternalErrorJSONResponse: api.InternalErrorJSONResponse{
				Error: stringPtr(err.Error()),
			},
		}, nil
	}
	return api.GetCommentById200JSONResponse(commentToResponse(comment)), nil
}

func (h *CommentHandler) UpdateComment(ctx context.Context, request api.UpdateCommentRequestObject) (api.UpdateCommentResponseObject, error) {
	err := h.commentService.Update(ctx, request.Id, request.Body.AuthorId, request.Body.PostId, request.Body.Text)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return api.UpdateComment404JSONResponse{
				NotFoundJSONResponse: api.NotFoundJSONResponse{
					Error: stringPtr("comment not found"),
				},
			}, nil
		}
		return api.UpdateComment500JSONResponse{
			InternalErrorJSONResponse: api.InternalErrorJSONResponse{
				Error: stringPtr(err.Error()),
			},
		}, nil
	}
	return api.UpdateComment204Response{}, nil
}

func (h *CommentHandler) DeleteComment(ctx context.Context, request api.DeleteCommentRequestObject) (api.DeleteCommentResponseObject, error) {
	err := h.commentService.Delete(ctx, request.Id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return api.DeleteComment404JSONResponse{
				NotFoundJSONResponse: api.NotFoundJSONResponse{
					Error: stringPtr("comment not found"),
				},
			}, nil
		}
		return api.DeleteComment500JSONResponse{
			InternalErrorJSONResponse: api.InternalErrorJSONResponse{
				Error: stringPtr(err.Error()),
			},
		}, nil
	}
	return api.DeleteComment204Response{}, nil
}
