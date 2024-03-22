package repositories

import (
	"context"
	"database/sql"

	"github.com/hutamatr/GoBlogify/exception"
	"github.com/hutamatr/GoBlogify/helpers"
	"github.com/hutamatr/GoBlogify/model/domain"
)

type CommentRepositoryImpl struct {
}

func NewCommentRepository() CommentRepository {
	return &CommentRepositoryImpl{}
}

func (repository *CommentRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, comment domain.Comment) domain.CommentJoin {
	ctxC, cancel := context.WithCancel(ctx)
	defer cancel()

	queryInsert := "INSERT INTO comment(post_id, user_id, content) VALUES(?, ?, ?)"

	result, err := tx.ExecContext(ctxC, queryInsert, comment.Post_Id, comment.User_Id, comment.Content)

	helpers.PanicError(err)

	id, err := result.LastInsertId()

	helpers.PanicError(err)

	createdComment := repository.FindById(ctx, tx, int(id))

	return createdComment
}

func (repository *CommentRepositoryImpl) FindCommentsByPost(ctx context.Context, tx *sql.Tx, postId int, offset int) []domain.CommentJoin {
	ctxC, cancel := context.WithCancel(ctx)
	defer cancel()

	query := "SELECT comment.id, comment.content, comment.post_id, comment.user_id, comment.created_at, comment.updated_at, user.id, user.username, user.email FROM comment INNER JOIN user ON user.id = comment.user_id WHERE post_id = ? LIMIT 10 OFFSET ?"

	rows, err := tx.QueryContext(ctxC, query, postId, offset)

	helpers.PanicError(err)

	defer rows.Close()

	var comments []domain.CommentJoin

	for rows.Next() {
		var comment domain.CommentJoin
		err := rows.Scan(&comment.Id, &comment.Content, &comment.Post_Id, &comment.User_Id, &comment.Created_At, &comment.Updated_At, &comment.User.Id, &comment.User.Username, &comment.User.Email)
		helpers.PanicError(err)

		comments = append(comments, comment)
	}

	return comments
}

func (repository *CommentRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, commentId int) domain.CommentJoin {
	ctxC, cancel := context.WithCancel(ctx)
	defer cancel()

	query := "SELECT comment.id, comment.content, comment.post_id, comment.user_id, comment.created_at, comment.updated_at, user.id, user.username, user.email FROM comment INNER JOIN user ON user.id = comment.user_id WHERE comment.id = ?"

	rows, err := tx.QueryContext(ctxC, query, commentId)

	helpers.PanicError(err)

	defer rows.Close()

	var comment domain.CommentJoin

	if rows.Next() {
		err := rows.Scan(&comment.Id, &comment.Content, &comment.Post_Id, &comment.User_Id, &comment.Created_At, &comment.Updated_At, &comment.User.Id, &comment.User.Username, &comment.User.Email)

		helpers.PanicError(err)
	} else {
		panic(exception.NewNotFoundError("comment not found"))
	}

	return comment
}

func (repository *CommentRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, comment domain.Comment) domain.CommentJoin {
	ctxC, cancel := context.WithCancel(ctx)
	defer cancel()

	query := "UPDATE comment SET content = ? WHERE id = ?"

	_, err := tx.ExecContext(ctxC, query, comment.Content, comment.Id)

	helpers.PanicError(err)

	updatedComment := repository.FindById(ctx, tx, comment.Id)

	return updatedComment
}

func (repository *CommentRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, commentId int) {
	ctxC, cancel := context.WithCancel(ctx)
	defer cancel()

	query := "DELETE FROM comment WHERE id = ?"

	result, err := tx.ExecContext(ctxC, query, commentId)

	helpers.PanicError(err)

	resultRows, err := result.RowsAffected()

	if resultRows == 0 {
		panic(exception.NewNotFoundError("comment not found"))
	}

	helpers.PanicError(err)
}
