package services

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/hutamatr/go-blog-api/helpers"
	"github.com/hutamatr/go-blog-api/model/domain"
	"github.com/hutamatr/go-blog-api/model/web"
	repositories "github.com/hutamatr/go-blog-api/repositories/articles"
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
	fmt.Println("errorr,", err)
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

func (service *ArticleServiceImpl) Update(ctx context.Context, article web.ArticleUpdateRequest) web.ArticleResponse {
	err := service.validator.Struct(article)
	helpers.PanicError(err)

	tx, err := service.db.Begin()
	helpers.PanicError(err)
	defer helpers.TxRollbackCommit(tx)

	service.repository.FindById(ctx, tx, article.Id)

	updateArticleData := domain.Article{
		Id:          article.Id,
		Title:       article.Title,
		Body:        article.Body,
		Author:      article.Author,
		Category_Id: article.Category_Id,
		Published:   article.Published,
		Deleted:     article.Deleted,
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
