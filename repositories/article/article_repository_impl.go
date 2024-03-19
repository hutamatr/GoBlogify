package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/hutamatr/go-blog-api/exception"
	"github.com/hutamatr/go-blog-api/helpers"
	"github.com/hutamatr/go-blog-api/model/domain"
)

type ArticleRepositoryImpl struct {
}

func NewArticleRepository() ArticleRepository {
	return &ArticleRepositoryImpl{}
}

func (repository *ArticleRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, article domain.Article) domain.ArticleJoinCategory {
	ctxC, cancel := context.WithCancel(ctx)
	defer cancel()

	queryInsert := "INSERT INTO article(title, body, author, published, category_id) VALUES(?, ?, ?, ?, ?)"

	result, err := tx.ExecContext(ctxC, queryInsert, article.Title, article.Body, article.Author, article.Published, article.Category_Id)

	helpers.PanicError(err)

	id, err := result.LastInsertId()

	helpers.PanicError(err)

	createdArticle := repository.FindById(ctx, tx, int(id))

	return createdArticle
}

func (repository *ArticleRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) []domain.ArticleJoinCategory {
	ctxC, cancel := context.WithCancel(ctx)
	defer cancel()

	query := "SELECT article.id, article.title, article.body, article.author, article.created_at, article.updated_at, article.deleted_at, article.deleted, article.published, category.id, category.name, category.created_at, category.updated_at FROM article INNER JOIN category ON article.category_id = category.id WHERE article.deleted = false"

	rows, err := tx.QueryContext(ctxC, query)

	helpers.PanicError(err)

	defer rows.Close()

	var articles []domain.ArticleJoinCategory

	var deletedAt sql.NullTime

	for rows.Next() {
		var article domain.ArticleJoinCategory

		err := rows.Scan(&article.Id, &article.Title, &article.Body, &article.Author, &article.Created_At, &article.Updated_At, &deletedAt, &article.Deleted, &article.Published, &article.Category.Id, &article.Category.Name, &article.Category.Created_At, &article.Category.Updated_At)

		helpers.PanicError(err)

		if deletedAt.Valid {
			article.Deleted_At = deletedAt.Time
		} else {
			article.Deleted_At = time.Time{}
		}

		articles = append(articles, article)
	}

	return articles
}

func (repository *ArticleRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, articleId int) domain.ArticleJoinCategory {
	ctxC, cancel := context.WithCancel(ctx)
	defer cancel()

	query := "SELECT article.id, article.title, article.body, article.author, article.created_at, article.updated_at, article.deleted_at, article.deleted, article.published, category.id, category.name, category.created_at, category.updated_at FROM article INNER JOIN category ON article.category_id = category.id WHERE article.id = ? AND article.deleted = false"

	rows, err := tx.QueryContext(ctxC, query, articleId)

	helpers.PanicError(err)

	defer rows.Close()

	var article domain.ArticleJoinCategory

	var deletedAt sql.NullTime

	if rows.Next() {
		err := rows.Scan(&article.Id, &article.Title, &article.Body, &article.Author, &article.Created_At, &article.Updated_At, &deletedAt, &article.Deleted, &article.Published, &article.Category.Id, &article.Category.Name, &article.Category.Created_At, &article.Category.Updated_At)

		helpers.PanicError(err)

		if deletedAt.Valid {
			article.Deleted_At = deletedAt.Time
		} else {
			article.Deleted_At = time.Time{}
		}
	} else {
		panic(exception.NewNotFoundError("article not found"))
	}

	return article
}

func (repository *ArticleRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, article domain.Article) domain.ArticleJoinCategory {
	ctxC, cancel := context.WithCancel(ctx)
	defer cancel()

	queryUpdate := "UPDATE article SET title = ?, body = ?, author = ?, category_id = ?, published = ?, deleted = ? WHERE id = ? AND deleted = false"

	_, err := tx.ExecContext(ctxC, queryUpdate, article.Title, article.Body, article.Author, article.Category_Id, article.Published, article.Deleted, article.Id)

	helpers.PanicError(err)

	updatedArticle := repository.FindById(ctx, tx, article.Id)

	return updatedArticle
}

func (repository *ArticleRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, articleId int) {
	ctxC, cancel := context.WithCancel(ctx)
	defer cancel()

	query := "UPDATE article SET deleted = true, deleted_at = NOW() WHERE id = ?"

	result, err := tx.ExecContext(ctxC, query, articleId)

	helpers.PanicError(err)

	resultRows, err := result.RowsAffected()

	if resultRows == 0 {
		panic(exception.NewNotFoundError("article not found"))
	}

	helpers.PanicError(err)
}
