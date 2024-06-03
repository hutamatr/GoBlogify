package post

import (
	"context"
	"database/sql"
	"fmt"
	"mime/multipart"

	"github.com/go-playground/validator/v10"
	"github.com/hutamatr/GoBlogify/exception"
	"github.com/hutamatr/GoBlogify/helpers"
	"github.com/hutamatr/GoBlogify/post_image"
)

type PostService interface {
	Create(ctx context.Context, request PostCreateRequest) PostResponse
	FindAllByUser(ctx context.Context, userId, limit, offset int) ([]PostResponse, int)
	FindAllByFollowed(ctx context.Context, userId, limit, offset int) ([]PostResponseFollowed, int)
	FindById(ctx context.Context, postId int) PostResponse
	Update(ctx context.Context, request PostUpdateRequest) PostResponse
	Delete(ctx context.Context, postId int)
}

type PostServiceImpl struct {
	postRepository      PostRepository
	postImageRepository post_image.PostImageRepository
	db                  *sql.DB
	validator           *validator.Validate
}

func NewPostService(postRepository PostRepository, postImageRepository post_image.PostImageRepository, db *sql.DB, validator *validator.Validate) PostService {
	return &PostServiceImpl{
		postRepository:      postRepository,
		postImageRepository: postImageRepository,
		db:                  db,
		validator:           validator,
	}
}

func (service *PostServiceImpl) Create(ctx context.Context, request PostCreateRequest) PostResponse {
	err := service.validator.Struct(request)
	helpers.PanicError(err, "invalid request")

	tx, err := service.db.Begin()
	helpers.PanicError(err, "failed to begin transaction")
	defer helpers.TxRollbackCommit(tx)

	imageRequest := []interface{}{request.Image_1, request.Image_2, request.Image_3}
	imageNameRequest := []interface{}{request.Image_Name_1, request.Image_Name_2, request.Image_Name_3}

	var imgResult []string
	for idx, img := range imageRequest {
		result, err := helpers.UploadToCloudinary(img, imageNameRequest[idx].(string))
		helpers.PanicError(err, "Failed upload image to cloud")
		imgResult = append(imgResult, result)
	}

	postRequest := Post{
		Title:       request.Title,
		Post_Body:   request.Post_Body,
		User_Id:     request.User_Id,
		Published:   request.Published,
		Category_Id: request.Category_Id,
	}

	createdPost := service.postRepository.Save(ctx, tx, postRequest)

	postImageRequest := post_image.PostImage{
		Post_Id:      createdPost.Id,
		Image_1:      imgResult[0],
		Image_Name_1: request.Image_Name_1,
		Image_2:      imgResult[1],
		Image_Name_2: request.Image_Name_2,
		Image_3:      imgResult[2],
		Image_Name_3: request.Image_Name_3,
	}

	createdPostImage := service.postImageRepository.Save(ctx, tx, postImageRequest)

	return ToPostResponse(createdPost, createdPostImage)
}

func (service *PostServiceImpl) FindAllByUser(ctx context.Context, userId, limit, offset int) ([]PostResponse, int) {
	tx, err := service.db.Begin()
	helpers.PanicError(err, "failed to begin transaction")
	defer helpers.TxRollbackCommit(tx)

	posts := service.postRepository.FindAllByUser(ctx, tx, userId, limit, offset)
	countPosts := service.postRepository.CountPostsByUser(ctx, tx, userId)
	var postsData []PostResponse

	if len(posts) == 0 {
		panic(exception.NewNotFoundError("posts not found"))
	}

	for _, post := range posts {
		postImageByPostId := service.postImageRepository.FindByPostId(ctx, tx, post.Id)
		postsData = append(postsData, ToPostResponse(post, postImageByPostId))
	}

	return postsData, countPosts
}

func (service *PostServiceImpl) FindAllByFollowed(ctx context.Context, userId, limit, offset int) ([]PostResponseFollowed, int) {
	tx, err := service.db.Begin()
	helpers.PanicError(err, "failed to begin transaction")
	defer helpers.TxRollbackCommit(tx)

	postsByFollowed := service.postRepository.FindAllByFollowed(ctx, tx, userId, limit, offset)

	var postByFollowedData []PostResponseFollowed

	if len(postsByFollowed) == 0 {
		panic(exception.NewNotFoundError("posts not found"))
	}

	for _, post := range postsByFollowed {
		postByFollowedData = append(postByFollowedData, ToPostResponseFollowed(post))
	}

	return postByFollowedData, len(postsByFollowed)
}

func (service *PostServiceImpl) FindById(ctx context.Context, postId int) PostResponse {
	tx, err := service.db.Begin()
	helpers.PanicError(err, "failed to begin transaction")
	defer helpers.TxRollbackCommit(tx)

	post := service.postRepository.FindById(ctx, tx, postId)
	postImageByPostId := service.postImageRepository.FindByPostId(ctx, tx, post.Id)

	return ToPostResponse(post, postImageByPostId)
}

func (service *PostServiceImpl) Update(ctx context.Context, request PostUpdateRequest) PostResponse {

	env := helpers.NewEnv()
	assetFolder := env.Cloudinary.AssetFolder

	err := service.validator.Struct(request)
	helpers.PanicError(err, "invalid request")

	tx, err := service.db.Begin()
	helpers.PanicError(err, "failed to begin transaction")
	defer helpers.TxRollbackCommit(tx)

	var updatePostData Post
	var imgResult []string
	var imageRequest []interface{}
	var imageNameRequest []string

	if attachment, exists := request.Image_1.(multipart.File); exists && request.Image_Name_1 != "" {
		imageRequest = append(imageRequest, attachment)
		imageNameRequest = append(imageNameRequest, request.Image_Name_1)
	}
	if attachment, exists := request.Image_2.(multipart.File); exists && request.Image_Name_2 != "" {
		imageRequest = append(imageRequest, attachment)
		imageNameRequest = append(imageNameRequest, request.Image_Name_2)
	}
	if attachment, exists := request.Image_3.(multipart.File); exists && request.Image_Name_3 != "" {
		imageRequest = append(imageRequest, attachment)
		imageNameRequest = append(imageNameRequest, request.Image_Name_3)
	}

	if len(imageRequest) > 0 {
		for idx, img := range imageRequest {
			result, err := helpers.UploadToCloudinary(img, imageNameRequest[idx])
			helpers.PanicError(err, "Failed upload image to cloud")
			imgResult = append(imgResult, result)
		}
	}

	updatePostData = Post{
		Id:          request.Id,
		Title:       request.Title,
		Post_Body:   request.Post_Body,
		User_Id:     request.User_Id,
		Category_Id: request.Category_Id,
		Published:   request.Published,
		Deleted:     request.Deleted,
	}

	service.postRepository.FindById(ctx, tx, request.Id)

	updatedPost := service.postRepository.Update(ctx, tx, updatePostData)

	postImage := service.postImageRepository.FindByPostId(ctx, tx, request.Id)

	var updatePostImageData post_image.PostImage

	updatePostImageData.Id = postImage.Id

	if request.Image_Name_1 != "" {
		updatePostImageData.Image_1 = imgResult[0]
		updatePostImageData.Image_Name_1 = request.Image_Name_1
		err := helpers.DeleteFromCloudinary(fmt.Sprintf("%s/%s", assetFolder, postImage.Image_Name_1))
		helpers.PanicError(err, "failed delete last image from cloud")
	}
	if request.Image_Name_2 != "" {
		updatePostImageData.Image_2 = imgResult[1]
		updatePostImageData.Image_Name_2 = request.Image_Name_2
		err := helpers.DeleteFromCloudinary(fmt.Sprintf("%s/%s", assetFolder, postImage.Image_Name_2))
		helpers.PanicError(err, "failed delete last image from cloud")
	}
	if request.Image_Name_3 != "" {
		updatePostImageData.Image_3 = imgResult[2]
		updatePostImageData.Image_Name_3 = request.Image_Name_3
		err := helpers.DeleteFromCloudinary(fmt.Sprintf("%s/%s", assetFolder, postImage.Image_Name_3))
		helpers.PanicError(err, "failed delete last image from cloud")
	}

	updatedPostImage := service.postImageRepository.Update(ctx, tx, updatePostImageData)

	helpers.PanicError(err, "failed to exec query update post")

	return ToPostResponse(updatedPost, updatedPostImage)
}

func (service *PostServiceImpl) Delete(ctx context.Context, postId int) {
	tx, err := service.db.Begin()
	helpers.PanicError(err, "failed to begin transaction")
	defer helpers.TxRollbackCommit(tx)

	service.postRepository.Delete(ctx, tx, postId)
}
