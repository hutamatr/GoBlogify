package post

import (
	"context"
	"database/sql"
	"time"

	"github.com/hutamatr/GoBlogify/exception"
	"github.com/hutamatr/GoBlogify/helpers"
)

type PostRepository interface {
	Save(ctx context.Context, tx *sql.Tx, post Post) PostJoin
	FindAll(ctx context.Context, tx *sql.Tx, limit, offset int) []PostJoin
	FindById(ctx context.Context, tx *sql.Tx, postId int) PostJoin
	Update(ctx context.Context, tx *sql.Tx, post Post) PostJoin
	Delete(ctx context.Context, tx *sql.Tx, postId int)
	CountPosts(ctx context.Context, tx *sql.Tx) int
}

type PostRepositoryImpl struct {
}

func NewPostRepository() PostRepository {
	return &PostRepositoryImpl{}
}

func (repository *PostRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, post Post) PostJoin {
	ctxC, cancel := context.WithCancel(ctx)
	defer cancel()

	queryInsert := "INSERT INTO post(title, body, is_published, user_id, category_id) VALUES(?, ?, ?, ?, ?)"

	result, err := tx.ExecContext(ctxC, queryInsert, post.Title, post.Body, post.Published, post.User_Id, post.Category_Id)

	helpers.PanicError(err, "failed to exec query insert post")

	id, err := result.LastInsertId()

	helpers.PanicError(err, "failed to get last insert id post")

	createdPost := repository.FindById(ctx, tx, int(id))

	return createdPost
}

func (repository *PostRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, limit, offset int) []PostJoin {
	query := "SELECT post.id, post.title, post.body, post.created_at, post.updated_at, post.deleted_at, post.is_deleted, post.is_published, user.id, user.role_id, user.username, user.email, user.first_name, user.last_name, user.created_at, user.updated_at, user.deleted_at, category.id, category.name, category.created_at, category.updated_at FROM post INNER JOIN category ON post.category_id = category.id INNER JOIN user ON post.user_id = user.id WHERE post.is_deleted = false LIMIT ? OFFSET ?"

	rows, err := tx.QueryContext(ctx, query, limit, offset)

	helpers.PanicError(err, "failed to query all posts")

	defer rows.Close()

	var posts []PostJoin

	var deletedAtPost sql.NullTime
	var deletedAtUser sql.NullTime
	var firstName sql.NullString
	var lastName sql.NullString

	for rows.Next() {
		var post PostJoin

		err := rows.Scan(&post.Id, &post.Title, &post.Body, &post.Created_At, &post.Updated_At, &deletedAtPost, &post.Deleted, &post.Published, &post.User.Id, &post.User.Role_Id, &post.User.Username, &post.User.Email, &firstName, &lastName, &post.User.Created_At, &post.User.Updated_At, &deletedAtUser, &post.Category.Id, &post.Category.Name, &post.Category.Created_At, &post.Category.Updated_At)

		helpers.PanicError(err, "failed to scan all posts")

		if deletedAtPost.Valid {
			post.Deleted_At = deletedAtPost.Time
		} else {
			post.Deleted_At = time.Time{}
		}
		if deletedAtUser.Valid {
			post.User.Deleted_At = deletedAtUser.Time
		} else {
			post.User.Deleted_At = time.Time{}
		}
		if firstName.Valid {
			post.User.First_Name = firstName.String
		} else {
			post.User.First_Name = ""
		}
		if lastName.Valid {
			post.User.Last_Name = lastName.String
		} else {
			post.User.Last_Name = ""
		}

		posts = append(posts, post)
	}

	return posts
}

func (repository *PostRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, postId int) PostJoin {
	query := "SELECT post.id, post.title, post.body, post.created_at, post.updated_at, post.deleted_at, post.is_deleted, post.is_published, user.id, user.role_id, user.username, user.email, user.first_name, user.last_name, user.created_at, user.updated_at, user.deleted_at, category.id, category.name, category.created_at, category.updated_at FROM post INNER JOIN category ON post.category_id = category.id INNER JOIN user ON post.user_id = user.id WHERE post.id = ? AND post.is_deleted = false"

	rows, err := tx.QueryContext(ctx, query, postId)

	helpers.PanicError(err, "failed to query post by id")

	defer rows.Close()

	var post PostJoin

	var deletedAtPost sql.NullTime
	var deletedAtUser sql.NullTime
	var firstName sql.NullString
	var lastName sql.NullString

	if rows.Next() {
		err := rows.Scan(&post.Id, &post.Title, &post.Body, &post.Created_At, &post.Updated_At, &deletedAtPost, &post.Deleted, &post.Published, &post.User.Id, &post.User.Role_Id, &post.User.Username, &post.User.Email, &firstName, &lastName, &post.User.Created_At, &post.User.Updated_At, &deletedAtUser, &post.Category.Id, &post.Category.Name, &post.Category.Created_At, &post.Category.Updated_At)

		helpers.PanicError(err, "failed to scan post by id")

		if deletedAtPost.Valid {
			post.Deleted_At = deletedAtPost.Time
		} else {
			post.Deleted_At = time.Time{}
		}
		if deletedAtUser.Valid {
			post.User.Deleted_At = deletedAtUser.Time
		} else {
			post.User.Deleted_At = time.Time{}
		}
		if firstName.Valid {
			post.User.First_Name = firstName.String
		} else {
			post.User.First_Name = ""
		}
		if lastName.Valid {
			post.User.Last_Name = lastName.String
		} else {
			post.User.Last_Name = ""
		}

	} else {
		panic(exception.NewNotFoundError("post not found"))
	}

	return post
}

func (repository *PostRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, post Post) PostJoin {
	queryUpdate := "UPDATE post SET title = ?, body = ?, category_id = ?, is_published = ?, is_deleted = ? WHERE id = ? AND is_deleted = false"

	_, err := tx.ExecContext(ctx, queryUpdate, post.Title, post.Body, post.Category_Id, post.Published, post.Deleted, post.Id)

	helpers.PanicError(err, "failed to exec query update post")

	updatedPost := repository.FindById(ctx, tx, post.Id)

	return updatedPost
}

func (repository *PostRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, postId int) {
	query := "UPDATE post SET is_deleted = true, deleted_at = NOW() WHERE id = ?"

	result, err := tx.ExecContext(ctx, query, postId)

	helpers.PanicError(err, "failed to exec query delete post")

	resultRows, err := result.RowsAffected()

	if resultRows == 0 {
		panic(exception.NewNotFoundError("post not found"))
	}

	helpers.PanicError(err, "failed to display rows affected delete post")
}

func (repository *PostRepositoryImpl) CountPosts(ctx context.Context, tx *sql.Tx) int {
	query := "SELECT COUNT(*) FROM post WHERE is_deleted = false"

	rows, err := tx.QueryContext(ctx, query)

	helpers.PanicError(err, "failed to query count posts")

	defer rows.Close()

	var countPosts int

	if rows.Next() {
		err := rows.Scan(&countPosts)
		helpers.PanicError(err, "failed to scan count posts")
	}

	return countPosts
}
