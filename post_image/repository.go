package post_image

import (
	"context"
	"database/sql"

	"github.com/hutamatr/GoBlogify/helpers"
)

type PostImageRepository interface {
	Save(ctx context.Context, tx *sql.Tx, postImage PostImage) PostImage
	FindById(ctx context.Context, tx *sql.Tx, postImageId int) PostImage
	FindByPostId(ctx context.Context, tx *sql.Tx, postId int) PostImage
	Update(ctx context.Context, tx *sql.Tx, postImage PostImage) PostImage
}

type PostImageRepositoryImpl struct{}

func NewPostImageRepository() PostImageRepository {
	return &PostImageRepositoryImpl{}
}

func (repository *PostImageRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, postImage PostImage) PostImage {
	query := "INSERT INTO post_images(post_id, image_1, image_name_1, image_2, image_name_2, image_3, image_name_3) VALUES (?, ?, ?, ?, ?, ?, ?)"

	result, err := tx.ExecContext(ctx, query, postImage.Post_Id, postImage.Image_1, postImage.Image_Name_1, postImage.Image_2, postImage.Image_Name_2, postImage.Image_3, postImage.Image_Name_3)

	helpers.PanicError(err, "failed to exec query insert post image")

	id, err := result.LastInsertId()

	helpers.PanicError(err, "failed to get last insert id post image")

	createdPostImage := repository.FindById(ctx, tx, int(id))

	return createdPostImage
}

func (repository *PostImageRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, postImageId int) PostImage {
	query := "SELECT id, post_id, image_1, image_name_1, image_2, image_name_2, image_3, image_name_3, created_at, updated_at FROM post_images WHERE id = ?"

	rows, err := tx.QueryContext(ctx, query, postImageId)

	helpers.PanicError(err, "failed to query post image by id")

	defer rows.Close()

	var postImages PostImage

	if rows.Next() {
		err := rows.Scan(&postImages.Id, &postImages.Post_Id, &postImages.Image_1, &postImages.Image_Name_1, &postImages.Image_2, &postImages.Image_Name_2, &postImages.Image_3, &postImages.Image_Name_3, &postImages.Created_At, &postImages.Updated_At)

		helpers.PanicError(err, "failed to scan post image by id")
	}

	return postImages
}

func (repository *PostImageRepositoryImpl) FindByPostId(ctx context.Context, tx *sql.Tx, postId int) PostImage {

	query := "SELECT id, image_1, image_name_1, image_2, image_name_2, image_3, image_name_3, created_at, updated_at FROM post_images WHERE post_id = ?"

	rows, err := tx.QueryContext(ctx, query, postId)

	helpers.PanicError(err, "failed to query post image by post id")

	defer rows.Close()

	var postImages PostImage

	if rows.Next() {
		err := rows.Scan(&postImages.Id, &postImages.Image_1, &postImages.Image_Name_1, &postImages.Image_2, &postImages.Image_Name_2, &postImages.Image_3, &postImages.Image_Name_3, &postImages.Created_At, &postImages.Updated_At)

		helpers.PanicError(err, "failed to scan post image by id")
	}

	return postImages
}

func (repository *PostImageRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, postImage PostImage) PostImage {
	var query string

	if postImage.Image_1 != "" && postImage.Image_Name_1 != "" {
		query = "UPDATE post_images SET image_1 = ?, image_name_1 = ? WHERE id = ?"

		_, err := tx.ExecContext(ctx, query, postImage.Image_1, postImage.Image_Name_1, postImage.Id)

		helpers.PanicError(err, "failed to exec query update post image_1")
	}
	if postImage.Image_2 != "" && postImage.Image_Name_2 != "" {
		query = "UPDATE post_images SET image_2 = ?, image_name_2 = ? WHERE id = ?"

		_, err := tx.ExecContext(ctx, query, postImage.Image_2, postImage.Image_Name_2, postImage.Id)

		helpers.PanicError(err, "failed to exec query update post image_2")
	}
	if postImage.Image_3 != "" && postImage.Image_Name_3 != "" {
		query = "UPDATE post_images SET image_3 = ?, image_name_3 = ? WHERE id = ?"

		_, err := tx.ExecContext(ctx, query, postImage.Image_3, postImage.Image_Name_3, postImage.Id)

		helpers.PanicError(err, "failed to exec query update post image_3")
	}

	updatedPost := repository.FindById(ctx, tx, postImage.Id)

	return updatedPost
}
