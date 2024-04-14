package post

import (
	"context"
	"database/sql"

	"github.com/go-playground/validator/v10"
	"github.com/hutamatr/GoBlogify/helpers"
)

type PostService interface {
	Create(ctx context.Context, request PostCreateRequest) PostResponse
	FindAll(ctx context.Context, limit, offset int) ([]PostResponse, int)
	FindById(ctx context.Context, postId int) PostResponse
	Update(ctx context.Context, request PostUpdateRequest) PostResponse
	Delete(ctx context.Context, postId int)
}

type PostServiceImpl struct {
	repository PostRepository
	db         *sql.DB
	validator  *validator.Validate
}

func NewPostService(postRepository PostRepository, db *sql.DB, validator *validator.Validate) PostService {
	return &PostServiceImpl{
		repository: postRepository,
		db:         db,
		validator:  validator,
	}
}

func (service *PostServiceImpl) Create(ctx context.Context, request PostCreateRequest) PostResponse {
	err := service.validator.Struct(request)
	helpers.PanicError(err, "invalid request")

	tx, err := service.db.Begin()
	helpers.PanicError(err, "failed to begin transaction")
	defer helpers.TxRollbackCommit(tx)

	postRequest := Post{
		Title:       request.Title,
		Body:        request.Body,
		User_Id:     request.User_Id,
		Published:   request.Published,
		Category_Id: request.Category_Id,
	}

	createdPost := service.repository.Save(ctx, tx, postRequest)

	return ToPostResponse(createdPost)
}

func (service *PostServiceImpl) FindAll(ctx context.Context, limit, offset int) ([]PostResponse, int) {
	tx, err := service.db.Begin()
	helpers.PanicError(err, "failed to begin transaction")
	defer helpers.TxRollbackCommit(tx)

	posts := service.repository.FindAll(ctx, tx, limit, offset)
	countPosts := service.repository.CountPosts(ctx, tx)

	var postsData []PostResponse

	if len(posts) == 0 {
		return postsData, 0
	}

	for _, post := range posts {
		postsData = append(postsData, ToPostResponse(post))
	}

	return postsData, countPosts
}

func (service *PostServiceImpl) FindById(ctx context.Context, postId int) PostResponse {
	tx, err := service.db.Begin()
	helpers.PanicError(err, "failed to begin transaction")
	defer helpers.TxRollbackCommit(tx)

	post := service.repository.FindById(ctx, tx, postId)

	return ToPostResponse(post)
}

func (service *PostServiceImpl) Update(ctx context.Context, request PostUpdateRequest) PostResponse {
	err := service.validator.Struct(request)
	helpers.PanicError(err, "invalid request")

	tx, err := service.db.Begin()
	helpers.PanicError(err, "failed to begin transaction")
	defer helpers.TxRollbackCommit(tx)

	service.repository.FindById(ctx, tx, request.Id)

	updatePostData := Post{
		Id:          request.Id,
		Title:       request.Title,
		Body:        request.Body,
		User_Id:     request.User_Id,
		Category_Id: request.Category_Id,
		Published:   request.Published,
		Deleted:     request.Deleted,
	}

	updatedPost := service.repository.Update(ctx, tx, updatePostData)

	helpers.PanicError(err, "failed to exec query update post")

	return ToPostResponse(updatedPost)
}

func (service *PostServiceImpl) Delete(ctx context.Context, postId int) {
	tx, err := service.db.Begin()
	helpers.PanicError(err, "failed to begin transaction")
	defer helpers.TxRollbackCommit(tx)

	service.repository.Delete(ctx, tx, postId)
}
