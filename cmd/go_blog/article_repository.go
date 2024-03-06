package go_blog

import (
	"context"
	"database/sql"
	"errors"

	"github.com/hutamatr/go-blog-api/cmd/go_blog/helper"
)

type ArticleRepository interface {
	Save(ctx context.Context, tx *sql.Tx, article Article) Article
	FindAll(ctx context.Context, tx *sql.Tx) []Article
	FindById(ctx context.Context, tx *sql.Tx, articleId int) (Article, error)
	Update(ctx context.Context, tx *sql.Tx, article Article) Article
	Delete(ctx context.Context, tx *sql.Tx, articleId int)
}

type ArticleRepositoryImpl struct {
}

func (repository *ArticleRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, article Article) Article {
	ctxC, cancel := context.WithCancel(ctx)
	defer cancel()

	query := "INSERT INTO article(title, body, author, published, category) VALUES(?, ?, ?, ?, ?)"

	result, err := tx.ExecContext(ctxC, query, article.Title, article.Body, article.Author, article.Published, article.Category)

	helper.PanicError(err)

	id, err := result.LastInsertId()

	helper.PanicError(err)

	article.Id = int(id)

	return article
}

func (repository *ArticleRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) []Article {
	ctxC, cancel := context.WithCancel(ctx)
	defer cancel()

	query := "SELECT id, title, body, author, created_at, updated_at, deleted_at, deleted, published, category FROM article"

	rows, err := tx.QueryContext(ctxC, query)

	helper.PanicError(err)

	defer rows.Close()

	articles := []Article{}

	for rows.Next() {
		article := Article{}

		err := rows.Scan(&article.Id, &article.Title, &article.Body, &article.Author, &article.Created_At, &article.Updated_At, &article.Deleted_At, &article.Deleted, &article.Published, &article.Category)

		helper.PanicError(err)

		articles = append(articles, article)
	}

	return articles
}

func (repository *ArticleRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, articleId int) (Article, error) {
	ctxC, cancel := context.WithCancel(ctx)
	defer cancel()

	query := "SELECT id, title, body, author, created_at, updated_at, deleted_at, deleted, published, category FROM article WHERE id = ? LIMIT = 1"

	rows, err := tx.QueryContext(ctxC, query, articleId)

	helper.PanicError(err)

	defer rows.Close()

	article := Article{}

	if rows.Next() {
		err := rows.Scan(&article.Id, &article.Title, &article.Body, &article.Author, &article.Created_At, &article.Updated_At, &article.Deleted_At, &article.Deleted, &article.Published, &article.Category)

		helper.PanicError(err)

		return article, nil
	} else {
		return article, errors.New("article not found")
	}
}

func (repository *ArticleRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, article Article) Article {
	ctxC, cancel := context.WithCancel(ctx)
	defer cancel()

	query := "UPDATE article SET title = ?, body = ?, author = ?, category = ?, published = ?, deleted = ? WHERE id = ?"

	result, err := tx.ExecContext(ctxC, query, article.Title, article.Body, article.Author, article.Category, article.Published, article.Deleted, article.Id)

	helper.PanicError(err)

	resultRows, err := result.RowsAffected()

	helper.PanicError(err)

	if resultRows == 0 {
		panic(errors.New("article not found"))
	}

	return article
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
		panic(errors.New("article not found"))
	}
}
