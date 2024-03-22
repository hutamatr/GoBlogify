package services

import (
	"context"
	"database/sql"

	"github.com/go-playground/validator/v10"
	"github.com/hutamatr/GoBlogify/helpers"
	"github.com/hutamatr/GoBlogify/model/domain"
	"github.com/hutamatr/GoBlogify/model/web"
	repositoriesPost "github.com/hutamatr/GoBlogify/repositories/post"
)

type PostServiceImpl struct {
	repository repositoriesPost.PostRepository
	db         *sql.DB
	validator  *validator.Validate
}

func NewPostService(postRepository repositoriesPost.PostRepository, db *sql.DB, validator *validator.Validate) PostService {
	return &PostServiceImpl{
		repository: postRepository,
		db:         db,
		validator:  validator,
	}
}

func (service *PostServiceImpl) Create(ctx context.Context, request web.PostCreateRequest) web.PostResponse {
	err := service.validator.Struct(request)
	helpers.PanicError(err)

	tx, err := service.db.Begin()
	helpers.PanicError(err)
	defer helpers.TxRollbackCommit(tx)

	postRequest := domain.Post{
		Title:       request.Title,
		Body:        request.Body,
		User_Id:     request.User_Id,
		Published:   request.Published,
		Category_Id: request.Category_Id,
	}

	createdPost := service.repository.Save(ctx, tx, postRequest)

	return web.ToPostResponse(createdPost)
}

func (service *PostServiceImpl) FindAll(ctx context.Context) []web.PostResponse {
	tx, err := service.db.Begin()
	helpers.PanicError(err)
	defer helpers.TxRollbackCommit(tx)

	posts := service.repository.FindAll(ctx, tx)

	var postsData []web.PostResponse

	if len(posts) == 0 {
		return postsData
	}

	for _, post := range posts {
		postsData = append(postsData, web.ToPostResponse(post))
	}

	return postsData
}

func (service *PostServiceImpl) FindById(ctx context.Context, postId int) web.PostResponse {
	tx, err := service.db.Begin()
	helpers.PanicError(err)
	defer helpers.TxRollbackCommit(tx)

	post := service.repository.FindById(ctx, tx, postId)

	return web.ToPostResponse(post)
}

func (service *PostServiceImpl) Update(ctx context.Context, request web.PostUpdateRequest) web.PostResponse {
	err := service.validator.Struct(request)
	helpers.PanicError(err)

	tx, err := service.db.Begin()
	helpers.PanicError(err)
	defer helpers.TxRollbackCommit(tx)

	service.repository.FindById(ctx, tx, request.Id)

	updatePostData := domain.Post{
		Id:          request.Id,
		Title:       request.Title,
		Body:        request.Body,
		User_Id:     request.User_Id,
		Category_Id: request.Category_Id,
		Published:   request.Published,
		Deleted:     request.Deleted,
	}

	updatedPost := service.repository.Update(ctx, tx, updatePostData)

	helpers.PanicError(err)

	return web.ToPostResponse(updatedPost)
}

func (service *PostServiceImpl) Delete(ctx context.Context, postId int) {
	tx, err := service.db.Begin()
	helpers.PanicError(err)
	defer helpers.TxRollbackCommit(tx)

	service.repository.Delete(ctx, tx, postId)
}
