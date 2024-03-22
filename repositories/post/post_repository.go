package repositories

import (
	"context"
	"database/sql"

	"github.com/hutamatr/GoBlogify/model/domain"
)

type PostRepository interface {
	Save(ctx context.Context, tx *sql.Tx, post domain.Post) domain.PostJoin
	FindAll(ctx context.Context, tx *sql.Tx) []domain.PostJoin
	FindById(ctx context.Context, tx *sql.Tx, postId int) domain.PostJoin
	Update(ctx context.Context, tx *sql.Tx, post domain.Post) domain.PostJoin
	Delete(ctx context.Context, tx *sql.Tx, postId int)
}
