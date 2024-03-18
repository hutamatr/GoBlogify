package services

import (
	"context"
	"database/sql"

	"github.com/go-playground/validator/v10"
	"github.com/hutamatr/go-blog-api/helpers"
	"github.com/hutamatr/go-blog-api/model/domain"
	"github.com/hutamatr/go-blog-api/model/web"
	repositories "github.com/hutamatr/go-blog-api/repositories/article"
)

type ArticleServiceImpl struct {
	repository repositories.ArticleRepository
	db         *sql.DB
	validator  *validator.Validate
}

func NewArticleService(articleRepository repositories.ArticleRepository, db *sql.DB, validator *validator.Validate) ArticleService {
	return &ArticleServiceImpl{
		repository: articleRepository,
		db:         db,
		validator:  validator,
	}
}

func (service *ArticleServiceImpl) Create(ctx context.Context, request web.ArticleCreateRequest) web.ArticleResponse {
	err := service.validator.Struct(request)
	helpers.PanicError(err)

	tx, err := service.db.Begin()
	helpers.PanicError(err)
	defer helpers.TxRollbackCommit(tx)

	articleRequest := domain.Article{
		Title:       request.Title,
		Body:        request.Body,
		Author:      request.Author,
		Published:   request.Published,
		Category_Id: request.Category_Id,
	}

	createdArticle := service.repository.Save(ctx, tx, articleRequest)

	return web.ArticleResponse(createdArticle)
}

func (service *ArticleServiceImpl) FindAll(ctx context.Context) []web.ArticleResponse {
	tx, err := service.db.Begin()
	helpers.PanicError(err)
	defer helpers.TxRollbackCommit(tx)

	articles := service.repository.FindAll(ctx, tx)

	var articlesData []web.ArticleResponse

	for _, article := range articles {
		articlesData = append(articlesData, web.ArticleResponse(article))
	}

	return articlesData
}

func (service *ArticleServiceImpl) FindById(ctx context.Context, articleId int) web.ArticleResponse {
	tx, err := service.db.Begin()
	helpers.PanicError(err)
	defer helpers.TxRollbackCommit(tx)

	article := service.repository.FindById(ctx, tx, articleId)

	return web.ArticleResponse(article)
}

func (service *ArticleServiceImpl) Update(ctx context.Context, request web.ArticleUpdateRequest) web.ArticleResponse {
	err := service.validator.Struct(request)
	helpers.PanicError(err)

	tx, err := service.db.Begin()
	helpers.PanicError(err)
	defer helpers.TxRollbackCommit(tx)

	service.repository.FindById(ctx, tx, request.Id)

	updateArticleData := domain.Article{
		Id:          request.Id,
		Title:       request.Title,
		Body:        request.Body,
		Author:      request.Author,
		Category_Id: request.Category_Id,
		Published:   request.Published,
		Deleted:     request.Deleted,
	}

	updatedArticle := service.repository.Update(ctx, tx, updateArticleData)

	helpers.PanicError(err)

	return web.ArticleResponse(updatedArticle)
}

func (service *ArticleServiceImpl) Delete(ctx context.Context, articleId int) {
	tx, err := service.db.Begin()
	helpers.PanicError(err)
	defer helpers.TxRollbackCommit(tx)

	service.repository.Delete(ctx, tx, articleId)
}
