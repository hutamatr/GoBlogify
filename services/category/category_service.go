package services

import (
	"context"

	"github.com/hutamatr/GoBlogify/model/web"
)

type CategoryService interface {
	Create(ctx context.Context, request web.CategoryCreateRequest, isAdmin bool) web.CategoryResponse
	FindAll(ctx context.Context) []web.CategoryResponse
	FindById(ctx context.Context, categoryId int) web.CategoryResponse
	Update(ctx context.Context, request web.CategoryUpdateRequest, isAdmin bool) web.CategoryResponse
	Delete(ctx context.Context, categoryId int, isAdmin bool)
}
