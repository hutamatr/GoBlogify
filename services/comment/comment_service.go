package services

import (
	"context"

	"github.com/hutamatr/GoBlogify/model/web"
)

type CommentService interface {
	Create(ctx context.Context, request web.CommentCreateRequest) web.CommentResponse
	FindCommentsByPost(ctx context.Context, postId, limit, offset int) ([]web.CommentResponse, int)
	FindById(ctx context.Context, commentId int) web.CommentResponse
	Update(ctx context.Context, request web.CommentUpdateRequest) web.CommentResponse
	Delete(ctx context.Context, commentId int)
}
