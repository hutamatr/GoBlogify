package services

import (
	"context"

	"github.com/hutamatr/GoBlogify/model/web"
)

type PostService interface {
	Create(ctx context.Context, request web.PostCreateRequest) web.PostResponse
	FindAll(ctx context.Context) []web.PostResponse
	FindById(ctx context.Context, postId int) web.PostResponse
	Update(ctx context.Context, request web.PostUpdateRequest) web.PostResponse
	Delete(ctx context.Context, postId int)
}
