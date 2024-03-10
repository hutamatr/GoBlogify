package go_blog

import (
	"context"
	"database/sql"

	"github.com/go-playground/validator/v10"
	"github.com/hutamatr/go-blog-api/cmd/go_blog/helper"
)

type ArticleService interface {
	Create(ctx context.Context, article ArticleCreateRequest) ArticleResponse
	FindAll(ctx context.Context) []ArticleResponse
	FindById(ctx context.Context, articleId int) ArticleResponse
	Update(ctx context.Context, article ArticleUpdateRequest) ArticleResponse
	Delete(ctx context.Context, articleId int)
}

type ArticleServiceImpl struct {
	repository ArticleRepository
	db         *sql.DB
	validator  *validator.Validate
}

func NewArticleService(articleRepository ArticleRepository, db *sql.DB, validator *validator.Validate) ArticleService {
	return &ArticleServiceImpl{
		repository: articleRepository,
		db:         db,
		validator:  validator,
	}
}

func (service *ArticleServiceImpl) Create(ctx context.Context, article ArticleCreateRequest) ArticleResponse {
	err := service.validator.Struct(article)
	helper.PanicError(err)

	tx, err := service.db.Begin()
	helper.PanicError(err)
	defer helper.TxRollbackCommit(tx)

	newArticle := ArticleCreateRequest{
		Title:       article.Title,
		Body:        article.Body,
		Author:      article.Author,
		Published:   article.Published,
		Category_Id: article.Category_Id,
	}

	createdArticle := service.repository.Save(ctx, tx, newArticle)

	return ArticleResponse(createdArticle)
}

func (service *ArticleServiceImpl) FindAll(ctx context.Context) []ArticleResponse {
	tx, err := service.db.Begin()
	helper.PanicError(err)
	defer helper.TxRollbackCommit(tx)

	articles := service.repository.FindAll(ctx, tx)

	var articlesData []ArticleResponse

	for _, article := range articles {
		articlesData = append(articlesData, ArticleResponse(article))
	}

	return articlesData
}

func (service *ArticleServiceImpl) FindById(ctx context.Context, articleId int) ArticleResponse {

	tx, err := service.db.Begin()
	helper.PanicError(err)
	defer helper.TxRollbackCommit(tx)

	article, err := service.repository.FindById(ctx, tx, articleId)
	helper.PanicError(err)

	return ArticleResponse(article)
}

func (service *ArticleServiceImpl) Update(ctx context.Context, article ArticleUpdateRequest) ArticleResponse {

	err := service.validator.Struct(article)
	helper.PanicError(err)

	tx, err := service.db.Begin()
	helper.PanicError(err)
	defer helper.TxRollbackCommit(tx)

	_, err = service.repository.FindById(ctx, tx, article.Id)
	helper.PanicError(err)

	updateArticleData := ArticleUpdateRequest{
		Id:          article.Id,
		Title:       article.Title,
		Body:        article.Body,
		Author:      article.Author,
		Category_Id: article.Category_Id,
		Published:   article.Published,
		Deleted:     article.Deleted,
	}

	// articleData.Title = article.Title
	// articleData.Body = article.Body
	// articleData.Author = article.Author
	// articleData.Category_Id = article.Category_Id
	// articleData.Published = article.Published
	// articleData.Deleted = article.Deleted

	updatedArticle := service.repository.Update(ctx, tx, updateArticleData)

	return updatedArticle
}

func (service *ArticleServiceImpl) Delete(ctx context.Context, articleId int) {
	tx, err := service.db.Begin()
	helper.PanicError(err)
	defer helper.TxRollbackCommit(tx)

	articleData, err := service.repository.FindById(ctx, tx, articleId)
	helper.PanicError(err)

	service.repository.Delete(ctx, tx, articleData.Id)
}
