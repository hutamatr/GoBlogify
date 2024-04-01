package repositories

import (
	"context"
	"database/sql"

	"github.com/hutamatr/GoBlogify/model/domain"
)

type CommentRepository interface {
	Save(ctx context.Context, tx *sql.Tx, comment domain.Comment) domain.CommentJoin
	FindCommentsByPost(ctx context.Context, tx *sql.Tx, postId, limit, offset int) []domain.CommentJoin
	FindById(ctx context.Context, tx *sql.Tx, commentId int) domain.CommentJoin
	Update(ctx context.Context, tx *sql.Tx, comment domain.Comment) domain.CommentJoin
	Delete(ctx context.Context, tx *sql.Tx, commentId int)
	CountCommentsByPost(ctx context.Context, tx *sql.Tx, postId int) int
}
