package go_blog

import (
	"context"
	"database/sql"
	"time"

	"github.com/hutamatr/go-blog-api/cmd/go_blog/helper"
)

type ArticleRepository interface {
	Save(ctx context.Context, tx *sql.Tx, article ArticleCreateRequest) ArticleResponse
	FindAll(ctx context.Context, tx *sql.Tx) []ArticleResponse
	FindById(ctx context.Context, tx *sql.Tx, articleId int) (ArticleResponse, error)
	Update(ctx context.Context, tx *sql.Tx, article ArticleUpdateRequest) ArticleResponse
	Delete(ctx context.Context, tx *sql.Tx, articleId int)
}

type ArticleRepositoryImpl struct {
}

func NewArticleRepository() ArticleRepository {
	return &ArticleRepositoryImpl{}
}

func (repository *ArticleRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, article ArticleCreateRequest) ArticleResponse {
	ctxC, cancel := context.WithCancel(ctx)
	defer cancel()

	queryInsert := "INSERT INTO article(title, body, author, published, category_id) VALUES(?, ?, ?, ?, ?)"

	result, err := tx.ExecContext(ctxC, queryInsert, article.Title, article.Body, article.Author, article.Published, article.Category_Id)

	helper.PanicError(err)

	id, err := result.LastInsertId()

	helper.PanicError(err)

	querySelect := "SELECT article.id, article.title, article.body, article.author, article.created_at, article.updated_at, article.deleted_at, article.deleted, article.published, category.id, category.name, category.created_at, category.updated_at FROM article INNER JOIN category ON article.category_id = category.id WHERE article.id = ?"

	rows, err := tx.QueryContext(ctxC, querySelect, id)

	helper.PanicError(err)

	defer rows.Close()

	createdArticle := ArticleResponse{}

	var deletedAt sql.NullTime

	if rows.Next() {
		err := rows.Scan(&createdArticle.Id, &createdArticle.Title, &createdArticle.Body, &createdArticle.Author, &createdArticle.Created_At, &createdArticle.Updated_At, &deletedAt, &createdArticle.Deleted, &createdArticle.Published, &createdArticle.Category.Id, &createdArticle.Category.Name, &createdArticle.Category.Created_At, &createdArticle.Category.Updated_At)

		helper.PanicError(err)

		if deletedAt.Valid {
			createdArticle.Deleted_At = deletedAt.Time
		} else {
			createdArticle.Deleted_At = time.Time{}
		}
	}

	return createdArticle
}

func (repository *ArticleRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) []ArticleResponse {
	ctxC, cancel := context.WithCancel(ctx)
	defer cancel()

	query := "SELECT article.id, article.title, article.body, article.author, article.created_at, article.updated_at, article.deleted_at, article.deleted, article.published, category.id, category.name, category.created_at, category.updated_at FROM article INNER JOIN category ON article.category_id = category.id"

	rows, err := tx.QueryContext(ctxC, query)

	helper.PanicError(err)

	defer rows.Close()

	articles := []ArticleResponse{}

	var deletedAt sql.NullTime

	for rows.Next() {
		article := ArticleResponse{}

		err := rows.Scan(&article.Id, &article.Title, &article.Body, &article.Author, &article.Created_At, &article.Updated_At, &deletedAt, &article.Deleted, &article.Published, &article.Category.Id, &article.Category.Name, &article.Category.Created_At, &article.Category.Updated_At)

		helper.PanicError(err)

		if deletedAt.Valid {
			article.Deleted_At = deletedAt.Time
		} else {
			article.Deleted_At = time.Time{}
		}

		articles = append(articles, article)
	}

	return articles
}

func (repository *ArticleRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, articleId int) (ArticleResponse, error) {
	ctxC, cancel := context.WithCancel(ctx)
	defer cancel()

	query := "SELECT article.id, article.title, article.body, article.author, article.created_at, article.updated_at, article.deleted_at, article.deleted, article.published, category.id, category.name, category.created_at, category.updated_at FROM article INNER JOIN category ON article.category_id = category.id WHERE article.id = ?"

	rows, err := tx.QueryContext(ctxC, query, articleId)

	helper.PanicError(err)

	defer rows.Close()

	var article ArticleResponse

	var deletedAt sql.NullTime

	if rows.Next() {
		err := rows.Scan(&article.Id, &article.Title, &article.Body, &article.Author, &article.Created_At, &article.Updated_At, &deletedAt, &article.Deleted, &article.Published, &article.Category.Id, &article.Category.Name, &article.Category.Created_At, &article.Category.Updated_At)

		helper.PanicError(err)

		if deletedAt.Valid {
			article.Deleted_At = deletedAt.Time
		} else {
			article.Deleted_At = time.Time{}
		}

		return article, nil
	} else {
		return article, helper.NotFoundError
	}
}

func (repository *ArticleRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, article ArticleUpdateRequest) ArticleResponse {
	ctxC, cancel := context.WithCancel(ctx)
	defer cancel()

	queryUpdate := "UPDATE article SET title = ?, body = ?, author = ?, category_id = ?, published = ?, deleted = ? WHERE id = ?"

	result, err := tx.ExecContext(ctxC, queryUpdate, article.Title, article.Body, article.Author, article.Category_Id, article.Published, article.Deleted, article.Id)

	helper.PanicError(err)

	resultRows, err := result.RowsAffected()

	helper.PanicError(err)

	if resultRows == 0 {
		panic(helper.NotFoundError)
	}

	querySelect := "SELECT article.id, article.title, article.body, article.author, article.created_at, article.updated_at, article.deleted_at, article.deleted, article.published, category.id, category.name, category.created_at, category.updated_at FROM article INNER JOIN category ON article.category_id = category.id WHERE article.id = ?"

	rows, err := tx.QueryContext(ctxC, querySelect, article.Id)

	helper.PanicError(err)

	defer rows.Close()

	var updatedArticle ArticleResponse

	var deletedAt sql.NullTime

	if rows.Next() {
		err := rows.Scan(&updatedArticle.Id, &updatedArticle.Title, &updatedArticle.Body, &updatedArticle.Author, &updatedArticle.Created_At, &updatedArticle.Updated_At, &deletedAt, &updatedArticle.Deleted, &updatedArticle.Published, &updatedArticle.Category.Id, &updatedArticle.Category.Name, &updatedArticle.Category.Created_At, &updatedArticle.Category.Updated_At)

		helper.PanicError(err)

		if deletedAt.Valid {
			updatedArticle.Deleted_At = deletedAt.Time
		} else {
			updatedArticle.Deleted_At = time.Time{}
		}
	}

	return updatedArticle
}

func (repository *ArticleRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, articleId int) {
	ctxC, cancel := context.WithCancel(ctx)
	defer cancel()

	query := "DELETE FROM article WHERE id = ?"

	result, err := tx.ExecContext(ctxC, query, articleId)

	helper.PanicError(err)

	resultRows, err := result.RowsAffected()

	helper.PanicError(err)

	if resultRows == 0 {
		panic(helper.NotFoundError)
	}
}
