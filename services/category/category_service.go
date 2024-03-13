package services

import (
	"context"

	"github.com/hutamatr/go-blog-api/model/web"
)

type CategoryService interface {
	Create(ctx context.Context, request web.CategoryCreateRequest) web.CategoryResponse
	FindAll(ctx context.Context) []web.CategoryResponse
	FindById(ctx context.Context, categoryId int) web.CategoryResponse
	Update(ctx context.Context, request web.CategoryUpdateRequest) web.CategoryResponse
	Delete(ctx context.Context, categoryId int)
}
