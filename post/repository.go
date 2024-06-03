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
	FindAllByFollowed(ctx context.Context, tx *sql.Tx, userId, limit, offset int) []PostJoinFollowed
	FindAllByUser(ctx context.Context, tx *sql.Tx, userId, limit, offset int) []PostJoin
	FindById(ctx context.Context, tx *sql.Tx, postId int) PostJoin
	Update(ctx context.Context, tx *sql.Tx, post Post) PostJoin
	Delete(ctx context.Context, tx *sql.Tx, postId int)
	CountPostsByUser(ctx context.Context, tx *sql.Tx, userId int) int
}

type PostRepositoryImpl struct {
}

func NewPostRepository() PostRepository {
	return &PostRepositoryImpl{}
}

func (repository *PostRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, post Post) PostJoin {
	queryInsert := "INSERT INTO posts(title, post_body, is_published, user_id, category_id) VALUES(?, ?, ?, ?, ?)"

	result, err := tx.ExecContext(ctx, queryInsert, post.Title, post.Post_Body, post.Published, post.User_Id, post.Category_Id)

	helpers.PanicError(err, "failed to exec query insert post")

	id, err := result.LastInsertId()

	helpers.PanicError(err, "failed to get last insert id post")

	createdPost := repository.FindById(ctx, tx, int(id))

	return createdPost
}

func (repository *PostRepositoryImpl) FindAllByUser(ctx context.Context, tx *sql.Tx, userId, limit, offset int) []PostJoin {

	query := `SELECT p.id, p.title, p.post_body, p.created_at, p.updated_at, p.deleted_at, p.is_deleted, p.is_published, u.id, u.role_id, u.username, u.email, u.first_name, u.last_name, u.created_at, u.updated_at, u.deleted_at, 
	(SELECT COUNT(*) FROM follows f WHERE f.followed_id = u.id) AS follower_count,
	(SELECT COUNT(*) FROM follows f WHERE f.follower_id = u.id) AS following_count,
	c.id, c.name, c.created_at, c.updated_at
	FROM users u 
	JOIN posts p 
	ON u.id = p.user_id 
	JOIN categories c 
	ON c.id = p.category_id 
	WHERE p.user_id = ? 
	AND p.is_deleted = false LIMIT ? OFFSET ?`

	rows, err := tx.QueryContext(ctx, query, userId, limit, offset)

	helpers.PanicError(err, "failed to query all posts")

	defer rows.Close()

	var posts []PostJoin

	var deletedAtPost sql.NullTime
	var deletedAtUser sql.NullTime
	var firstName sql.NullString
	var lastName sql.NullString

	for rows.Next() {
		var post PostJoin

		err := rows.Scan(&post.Id, &post.Title, &post.Post_Body, &post.Created_At, &post.Updated_At, &deletedAtPost, &post.Deleted, &post.Published, &post.User.Id, &post.User.Role_Id, &post.User.Username, &post.User.Email, &firstName, &lastName, &post.User.Created_At, &post.User.Updated_At, &deletedAtUser, &post.User.Follower, &post.User.Following, &post.Category.Id, &post.Category.Name, &post.Category.Created_At, &post.Category.Updated_At)

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

func (repository *PostRepositoryImpl) FindAllByFollowed(ctx context.Context, tx *sql.Tx, userId, limit, offset int) []PostJoinFollowed {

	query := `SELECT p.id, p.title, p.post_body, p.created_at, p.updated_at, p.deleted_at, p.is_deleted, p.is_published, u.id, u.role_id, u.username, u.email, u.first_name, u.last_name, u.created_at, u.updated_at, u.deleted_at, 
	(SELECT COUNT(*) FROM follows f WHERE f.followed_id = u.id) AS follower_count,
	(SELECT COUNT(*) FROM follows f WHERE f.follower_id = u.id) AS following_count
	FROM users u 
	JOIN posts p 
	ON u.id = p.user_id 
	JOIN follows f 
	ON u.id = f.followed_id
	WHERE f.follower_id = ? 
	AND p.is_deleted = false 
	ORDER BY p.created_at DESC LIMIT ? OFFSET ?`

	rows, err := tx.QueryContext(ctx, query, userId, limit, offset)

	helpers.PanicError(err, "failed to query post by user followed")

	defer rows.Close()

	var postsByFollowed []PostJoinFollowed

	var deletedAtPost sql.NullTime
	var deletedAtUser sql.NullTime
	var firstName sql.NullString
	var lastName sql.NullString

	for rows.Next() {
		var postByFollowed PostJoinFollowed
		err := rows.Scan(&postByFollowed.Id, &postByFollowed.Title, &postByFollowed.Post_Body, &postByFollowed.Created_At, &postByFollowed.Updated_At, &deletedAtPost, &postByFollowed.Deleted, &postByFollowed.Published, &postByFollowed.User.Id, &postByFollowed.User.Role_Id, &postByFollowed.User.Username, &postByFollowed.User.Email, &firstName, &lastName, &postByFollowed.User.Created_At, &postByFollowed.User.Updated_At, &deletedAtUser, &postByFollowed.User.Follower, &postByFollowed.User.Following)

		helpers.PanicError(err, "failed to scan post by user followed")

		if deletedAtPost.Valid {
			postByFollowed.Deleted_At = deletedAtPost.Time
		} else {
			postByFollowed.Deleted_At = time.Time{}
		}
		if deletedAtUser.Valid {
			postByFollowed.User.Deleted_At = deletedAtUser.Time
		} else {
			postByFollowed.User.Deleted_At = time.Time{}
		}

		if firstName.Valid {
			postByFollowed.User.First_Name = firstName.String
		} else {
			postByFollowed.User.First_Name = ""
		}
		if lastName.Valid {
			postByFollowed.User.Last_Name = lastName.String
		} else {
			postByFollowed.User.Last_Name = ""
		}

		postsByFollowed = append(postsByFollowed, postByFollowed)
	}

	return postsByFollowed
}

func (repository *PostRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, postId int) PostJoin {
	query := `SELECT p.id, p.title, p.post_body, p.created_at, p.updated_at, p.deleted_at, p.is_deleted, p.is_published, u.id, u.role_id, u.username, u.email, u.first_name, u.last_name, u.created_at, u.updated_at, u.deleted_at, c.id, c.name, c.created_at, c.updated_at 
	FROM users u 
	JOIN posts p 
	ON u.id = p.user_id 
	JOIN categories c 
	ON p.category_id = c.id  
	WHERE p.id = ? AND p.is_deleted = false`

	rows, err := tx.QueryContext(ctx, query, postId)

	helpers.PanicError(err, "failed to query post by id")

	defer rows.Close()

	var post PostJoin

	var deletedAtPost sql.NullTime
	var deletedAtUser sql.NullTime
	var firstName sql.NullString
	var lastName sql.NullString

	if rows.Next() {
		err := rows.Scan(&post.Id, &post.Title, &post.Post_Body, &post.Created_At, &post.Updated_At, &deletedAtPost, &post.Deleted, &post.Published, &post.User.Id, &post.User.Role_Id, &post.User.Username, &post.User.Email, &firstName, &lastName, &post.User.Created_At, &post.User.Updated_At, &deletedAtUser, &post.Category.Id, &post.Category.Name, &post.Category.Created_At, &post.Category.Updated_At)

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
	query := "UPDATE posts SET title = ?, post_body = ?, category_id = ?, is_published = ?, is_deleted = ? WHERE id = ? AND is_deleted = false"

	_, err := tx.ExecContext(ctx, query, post.Title, post.Post_Body, post.Category_Id, post.Published, post.Deleted, post.Id)

	helpers.PanicError(err, "failed to exec query update post")

	updatedPost := repository.FindById(ctx, tx, post.Id)

	return updatedPost
}

func (repository *PostRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, postId int) {
	query := "UPDATE posts SET is_deleted = true, deleted_at = NOW() WHERE id = ?"

	result, err := tx.ExecContext(ctx, query, postId)

	helpers.PanicError(err, "failed to exec query delete post")

	resultRows, err := result.RowsAffected()

	if resultRows == 0 {
		panic(exception.NewNotFoundError("post not found"))
	}

	helpers.PanicError(err, "failed to display rows affected delete post")
}

func (repository *PostRepositoryImpl) CountPostsByUser(ctx context.Context, tx *sql.Tx, userId int) int {
	query := "SELECT COUNT(*) FROM posts WHERE is_deleted = false AND user_id = ?"

	rows, err := tx.QueryContext(ctx, query, userId)

	helpers.PanicError(err, "failed to query count posts")

	defer rows.Close()

	var countPosts int

	if rows.Next() {
		err := rows.Scan(&countPosts)
		helpers.PanicError(err, "failed to scan count posts")
	}

	return countPosts
}
