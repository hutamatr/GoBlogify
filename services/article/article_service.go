package services

import (
	"context"

	"github.com/hutamatr/go-blog-api/model/web"
)

type ArticleService interface {
	Create(ctx context.Context, request web.ArticleCreateRequest) web.ArticleResponse
	FindAll(ctx context.Context) []web.ArticleResponse
	FindById(ctx context.Context, articleId int) web.ArticleResponse
	Update(ctx context.Context, request web.ArticleUpdateRequest) web.ArticleResponse
	Delete(ctx context.Context, articleId int)
}
