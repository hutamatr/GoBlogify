package go_blog

import "context"

type CategoryService interface {
	Create(ctx context.Context, article CategoryCreateRequest) CategoryResponse
	FindAll(ctx context.Context) []CategoryResponse
	FindById(ctx context.Context, articleId int) CategoryResponse
	Update(ctx context.Context, articleId int) CategoryResponse
	Delete(ctx context.Context, articleId int)
}
