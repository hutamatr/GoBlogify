package repositories

import (
	"context"
	"database/sql"

	"github.com/hutamatr/go-blog-api/model/domain"
)

type ArticleRepository interface {
	Save(ctx context.Context, tx *sql.Tx, article domain.Article) domain.ArticleJoinCategory
	FindAll(ctx context.Context, tx *sql.Tx) []domain.ArticleJoinCategory
	FindById(ctx context.Context, tx *sql.Tx, articleId int) domain.ArticleJoinCategory
	Update(ctx context.Context, tx *sql.Tx, article domain.Article) domain.ArticleJoinCategory
	Delete(ctx context.Context, tx *sql.Tx, articleId int)
}
