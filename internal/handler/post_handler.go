package handler

import (
	"context"
	"errors"

	"crud_service/api/api"
	"crud_service/internal/repository"
	"crud_service/internal/service"
)

type PostHandler struct {
	postService service.PostService
}

func NewPostHandler(postService service.PostService) *PostHandler {
	return &PostHandler{postService: postService}
}

func (h *PostHandler) CreatePost(ctx context.Context, request api.CreatePostRequestObject) (api.CreatePostResponseObject, error) {
	post, err := h.postService.Create(ctx, request.Body.Name, request.Body.Text, request.Body.AuthorId)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return api.CreatePost404JSONResponse{
				NotFoundJSONResponse: api.NotFoundJSONResponse{
					Error: stringPtr("author not found"),
				},
			}, nil
		}
		return api.CreatePost500JSONResponse{
			InternalErrorJSONResponse: api.InternalErrorJSONResponse{
				Error: stringPtr(err.Error()),
			},
		}, nil
	}
	return api.CreatePost201JSONResponse(postToResponse(post)), nil
}

func (h *PostHandler) ListPosts(ctx context.Context, request api.ListPostsRequestObject) (api.ListPostsResponseObject, error) {
	posts, err := h.postService.GetAll(ctx)
	if err != nil {
		return api.ListPosts500JSONResponse{
			InternalErrorJSONResponse: api.InternalErrorJSONResponse{
				Error: stringPtr(err.Error()),
			},
		}, nil
	}

	resp := make([]api.PostResponse, 0, len(posts))
	for _, p := range posts {
		resp = append(resp, postToResponse(p))
	}
	return api.ListPosts200JSONResponse(resp), nil
}

func (h *PostHandler) GetPostById(ctx context.Context, request api.GetPostByIdRequestObject) (api.GetPostByIdResponseObject, error) {
	post, err := h.postService.GetById(ctx, request.Id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return api.GetPostById404JSONResponse{
				NotFoundJSONResponse: api.NotFoundJSONResponse{
					Error: stringPtr("post not found"),
				},
			}, nil
		}
		return api.GetPostById500JSONResponse{
			InternalErrorJSONResponse: api.InternalErrorJSONResponse{
				Error: stringPtr(err.Error()),
			},
		}, nil
	}
	return api.GetPostById200JSONResponse(postToResponse(post)), nil
}

func (h *PostHandler) UpdatePost(ctx context.Context, request api.UpdatePostRequestObject) (api.UpdatePostResponseObject, error) {
	err := h.postService.Update(ctx, request.Body.Name, request.Body.Text, request.Id, request.Body.AuthorId)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return api.UpdatePost404JSONResponse{
				NotFoundJSONResponse: api.NotFoundJSONResponse{
					Error: stringPtr("post not found"),
				},
			}, nil
		}
		return api.UpdatePost500JSONResponse{
			InternalErrorJSONResponse: api.InternalErrorJSONResponse{
				Error: stringPtr(err.Error()),
			},
		}, nil
	}
	return api.UpdatePost204Response{}, nil
}

func (h *PostHandler) DeletePost(ctx context.Context, request api.DeletePostRequestObject) (api.DeletePostResponseObject, error) {
	err := h.postService.Delete(ctx, request.Id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return api.DeletePost404JSONResponse{
				NotFoundJSONResponse: api.NotFoundJSONResponse{
					Error: stringPtr("post not found"),
				},
			}, nil
		}
		return api.DeletePost500JSONResponse{
			InternalErrorJSONResponse: api.InternalErrorJSONResponse{
				Error: stringPtr(err.Error()),
			},
		}, nil
	}
	return api.DeletePost204Response{}, nil
}
