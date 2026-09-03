package handler

import (
	"crud_service/api/api"
	"crud_service/internal/service"
)

type Server struct {
	*UserHandler
	*PostHandler
	*CommentHandler
}

func NewServer(
	userService service.UserService,
	postService service.PostService,
	commentService service.CommentService,
) api.StrictServerInterface {
	return &Server{
		UserHandler:    NewUserHandler(userService),
		PostHandler:    NewPostHandler(postService),
		CommentHandler: NewCommentHandler(commentService),
	}
}
