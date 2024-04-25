package comment

import (
	"context"
	"database/sql"

	"github.com/hutamatr/GoBlogify/exception"
	"github.com/hutamatr/GoBlogify/helpers"
)

type CommentRepository interface {
	Save(ctx context.Context, tx *sql.Tx, comment Comment) CommentJoin
	FindCommentsByPost(ctx context.Context, tx *sql.Tx, postId, limit, offset int) []CommentJoin
	FindById(ctx context.Context, tx *sql.Tx, commentId int) CommentJoin
	Update(ctx context.Context, tx *sql.Tx, comment Comment) CommentJoin
	Delete(ctx context.Context, tx *sql.Tx, commentId int)
	CountCommentsByPost(ctx context.Context, tx *sql.Tx, postId int) int
}

type CommentRepositoryImpl struct {
}

func NewCommentRepository() CommentRepository {
	return &CommentRepositoryImpl{}
}

func (repository *CommentRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, comment Comment) CommentJoin {
	queryInsert := "INSERT INTO comment(post_id, user_id, content) VALUES(?, ?, ?)"

	result, err := tx.ExecContext(ctx, queryInsert, comment.Post_Id, comment.User_Id, comment.Content)

	helpers.PanicError(err, "failed to exec query insert comment")

	id, err := result.LastInsertId()

	helpers.PanicError(err, "failed to get last insert id comment")

	createdComment := repository.FindById(ctx, tx, int(id))

	return createdComment
}

func (repository *CommentRepositoryImpl) FindCommentsByPost(ctx context.Context, tx *sql.Tx, postId, limit, offset int) []CommentJoin {
	query := `SELECT c.id, c.content, c.post_id, c.user_id, c.created_at, c.updated_at, u.id, u.username, u.email 
	FROM user u 
	JOIN comment c 
	ON u.id = c.user_id 
	WHERE c.post_id = ? LIMIT ? OFFSET ?`

	rows, err := tx.QueryContext(ctx, query, postId, limit, offset)

	helpers.PanicError(err, "failed to query comments by post")

	defer rows.Close()

	var comments []CommentJoin

	for rows.Next() {
		var comment CommentJoin
		err := rows.Scan(&comment.Id, &comment.Content, &comment.Post_Id, &comment.User_Id, &comment.Created_At, &comment.Updated_At, &comment.User.Id, &comment.User.Username, &comment.User.Email)
		helpers.PanicError(err, "failed to scan comments by post")

		comments = append(comments, comment)
	}

	return comments
}

func (repository *CommentRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, commentId int) CommentJoin {
	query := "SELECT c.id, c.content, c.post_id, c.user_id, c.created_at, c.updated_at, u.id, u.username, u.email FROM user u JOIN comment c ON u.id = c.user_id WHERE c.id = ?"

	rows, err := tx.QueryContext(ctx, query, commentId)

	helpers.PanicError(err, "failed to query comment by id")

	defer rows.Close()

	var comment CommentJoin

	if rows.Next() {
		err := rows.Scan(&comment.Id, &comment.Content, &comment.Post_Id, &comment.User_Id, &comment.Created_At, &comment.Updated_At, &comment.User.Id, &comment.User.Username, &comment.User.Email)

		helpers.PanicError(err, "failed to scan comment by id")
	} else {
		panic(exception.NewNotFoundError("comment not found"))
	}

	return comment
}

func (repository *CommentRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, comment Comment) CommentJoin {
	query := "UPDATE comment SET content = ? WHERE id = ?"

	_, err := tx.ExecContext(ctx, query, comment.Content, comment.Id)

	helpers.PanicError(err, "failed to exec query update comment")

	updatedComment := repository.FindById(ctx, tx, comment.Id)

	return updatedComment
}

func (repository *CommentRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, commentId int) {
	query := "DELETE FROM comment WHERE id = ?"

	result, err := tx.ExecContext(ctx, query, commentId)

	helpers.PanicError(err, "failed to exec query delete comment")

	resultRows, err := result.RowsAffected()

	if resultRows == 0 {
		panic(exception.NewNotFoundError("comment not found"))
	}

	helpers.PanicError(err, "failed to display rows affected delete comment")
}

func (repository *CommentRepositoryImpl) CountCommentsByPost(ctx context.Context, tx *sql.Tx, postId int) int {
	query := "SELECT COUNT(*) FROM comment WHERE post_id = ?"

	rows, err := tx.QueryContext(ctx, query, postId)

	helpers.PanicError(err, "failed to query count comments by post")

	defer rows.Close()

	var countComments int

	if rows.Next() {
		err := rows.Scan(&countComments)
		helpers.PanicError(err, "failed to scan count comments by post")
	}

	return countComments
}
