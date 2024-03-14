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

	querySelect := "SELECT article.id, article.title, article.body, article.author, article.created_at, article.updated_at, article.deleted_at, article.deleted, article.published, category.id, category.name, category.created_at, category.updated_at FROM article INNER JOIN category ON article.category_id = category.id WHERE article.id = ?"

	rows, err := tx.QueryContext(ctxC, querySelect, id)

	helpers.PanicError(err)

	defer rows.Close()

	var createdArticle domain.ArticleJoinCategory

	var deletedAt sql.NullTime

	if rows.Next() {
		err := rows.Scan(&createdArticle.Id, &createdArticle.Title, &createdArticle.Body, &createdArticle.Author, &createdArticle.Created_At, &createdArticle.Updated_At, &deletedAt, &createdArticle.Deleted, &createdArticle.Published, &createdArticle.Category.Id, &createdArticle.Category.Name, &createdArticle.Category.Created_At, &createdArticle.Category.Updated_At)

		helpers.PanicError(err)

		if deletedAt.Valid {
			createdArticle.Deleted_At = deletedAt.Time
		} else {
			createdArticle.Deleted_At = time.Time{}
		}
	}

	return createdArticle
}

func (repository *ArticleRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) []domain.ArticleJoinCategory {
	ctxC, cancel := context.WithCancel(ctx)
	defer cancel()

	query := "SELECT article.id, article.title, article.body, article.author, article.created_at, article.updated_at, article.deleted_at, article.deleted, article.published, category.id, category.name, category.created_at, category.updated_at FROM article INNER JOIN category ON article.category_id = category.id"

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

	query := "SELECT article.id, article.title, article.body, article.author, article.created_at, article.updated_at, article.deleted_at, article.deleted, article.published, category.id, category.name, category.created_at, category.updated_at FROM article INNER JOIN category ON article.category_id = category.id WHERE article.id = ?"

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

	queryUpdate := "UPDATE article SET title = ?, body = ?, author = ?, category_id = ?, published = ?, deleted = ? WHERE id = ?"

	_, err := tx.ExecContext(ctxC, queryUpdate, article.Title, article.Body, article.Author, article.Category_Id, article.Published, article.Deleted, article.Id)

	helpers.PanicError(err)

	querySelect := "SELECT article.id, article.title, article.body, article.author, article.created_at, article.updated_at, article.deleted_at, article.deleted, article.published, category.id, category.name, category.created_at, category.updated_at FROM article INNER JOIN category ON article.category_id = category.id WHERE article.id = ?"

	rows, err := tx.QueryContext(ctxC, querySelect, article.Id)

	helpers.PanicError(err)

	defer rows.Close()

	var updatedArticle domain.ArticleJoinCategory

	var deletedAt sql.NullTime

	if rows.Next() {
		err := rows.Scan(&updatedArticle.Id, &updatedArticle.Title, &updatedArticle.Body, &updatedArticle.Author, &updatedArticle.Created_At, &updatedArticle.Updated_At, &deletedAt, &updatedArticle.Deleted, &updatedArticle.Published, &updatedArticle.Category.Id, &updatedArticle.Category.Name, &updatedArticle.Category.Created_At, &updatedArticle.Category.Updated_At)

		helpers.PanicError(err)

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

	helpers.PanicError(err)

	resultRows, err := result.RowsAffected()

	if resultRows == 0 {
		panic(exception.NewNotFoundError("article not found"))
	}

	helpers.PanicError(err)
}
