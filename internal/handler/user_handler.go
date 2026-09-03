package handler

import (
	"context"
	"crud_service/api/api"
	"crud_service/internal/repository"
	"crud_service/internal/service"
	"errors"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) CreateUser(ctx context.Context, request api.CreateUserRequestObject) (api.CreateUserResponseObject, error) {
	name := request.Body.Name
	email := string(request.Body.Email)

	if name == "" || email == "" {
		return api.CreateUser400JSONResponse{
			BadRequestJSONResponse: api.BadRequestJSONResponse{
				Error: stringPtr("give me more data, bitch"),
			},
		}, nil
	}

	user, err := h.userService.Create(
		ctx,
		name,
		email,
	)

	if err != nil {
		return api.CreateUser500JSONResponse{
			InternalErrorJSONResponse: api.InternalErrorJSONResponse{
				Error: stringPtr(err.Error()),
			},
		}, nil
	}

	return api.CreateUser201JSONResponse(userToResponse(user)), nil
}

func (h *UserHandler) ListUsers(ctx context.Context, request api.ListUsersRequestObject) (api.ListUsersResponseObject, error) {
	users, err := h.userService.GetAll(ctx)
	if err != nil {
		return api.ListUsers500JSONResponse{
			InternalErrorJSONResponse: api.InternalErrorJSONResponse{
				Error: stringPtr(err.Error()),
			},
		}, nil
	}

	users_response := make([]api.UserResponse, 0, len(users))

	for _, user := range users {
		users_response = append(users_response, userToResponse(user))
	}

	return api.ListUsers200JSONResponse(users_response), nil
}

func (h *UserHandler) GetUserById(ctx context.Context, request api.GetUserByIdRequestObject) (api.GetUserByIdResponseObject, error) {
	id := request.Id

	user, err := h.userService.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return api.GetUserById404JSONResponse{
				NotFoundJSONResponse: api.NotFoundJSONResponse{
					Error: stringPtr(err.Error()),
				},
			}, nil
		}
		return api.GetUserById500JSONResponse{
			InternalErrorJSONResponse: api.InternalErrorJSONResponse{
				Error: stringPtr(err.Error()),
			},
		}, nil
	}

	return api.GetUserById200JSONResponse(userToResponse(user)), nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, request api.UpdateUserRequestObject) (api.UpdateUserResponseObject, error) {
	name := request.Body.Name
	email := string(request.Body.Email)
	id := request.Id

	if name == "" || email == "" || id < 0 {
		return api.UpdateUser400JSONResponse{
			BadRequestJSONResponse: api.BadRequestJSONResponse{
				Error: stringPtr("bad request"),
			},
		}, nil
	}

	err := h.userService.Update(ctx, name, email, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return api.UpdateUser404JSONResponse{
				NotFoundJSONResponse: api.NotFoundJSONResponse{
					Error: stringPtr(err.Error()),
				},
			}, nil
		}
		return api.UpdateUser500JSONResponse{
			InternalErrorJSONResponse: api.InternalErrorJSONResponse{
				Error: stringPtr(err.Error()),
			},
		}, nil
	}

	return api.UpdateUser204Response{}, nil
}

func (h *UserHandler) DeleteUser(ctx context.Context, request api.DeleteUserRequestObject) (api.DeleteUserResponseObject, error) {
	err := h.userService.Delete(ctx, request.Id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return api.DeleteUser404JSONResponse{
				NotFoundJSONResponse: api.NotFoundJSONResponse{
					Error: stringPtr("user not found"),
				},
			}, nil
		}
		return api.DeleteUser500JSONResponse{
			InternalErrorJSONResponse: api.InternalErrorJSONResponse{
				Error: stringPtr(err.Error()),
			},
		}, nil
	}
	return api.DeleteUser204Response{}, nil
}
