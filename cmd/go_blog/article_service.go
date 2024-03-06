package go_blog

import "context"

type ArticleService interface {
	Create(ctx context.Context, article ArticleCreateRequest) ArticleResponse
	FindAll(ctx context.Context) []ArticleResponse
	FindById(ctx context.Context, articleId int) ArticleResponse
	Update(ctx context.Context, articleId int) ArticleResponse
	Delete(ctx context.Context, articleId int)
}
