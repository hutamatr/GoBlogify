package services

import (
	"context"
	"database/sql"

	"github.com/go-playground/validator/v10"
	"github.com/hutamatr/GoBlogify/helpers"
	"github.com/hutamatr/GoBlogify/model/domain"
	"github.com/hutamatr/GoBlogify/model/web"
	repositoriesComment "github.com/hutamatr/GoBlogify/repositories/comment"
)

type CommentServiceImpl struct {
	repository repositoriesComment.CommentRepository
	db         *sql.DB
	validator  *validator.Validate
}

func NewCommentService(commentRepository repositoriesComment.CommentRepository, db *sql.DB, validator *validator.Validate) CommentService {
	return &CommentServiceImpl{
		repository: commentRepository,
		db:         db,
		validator:  validator,
	}
}

func (service *CommentServiceImpl) Create(ctx context.Context, request web.CommentCreateRequest) web.CommentResponse {
	err := service.validator.Struct(request)
	helpers.PanicError(err)

	tx, err := service.db.Begin()
	helpers.PanicError(err)
	defer helpers.TxRollbackCommit(tx)

	newComment := domain.Comment{
		Post_Id: request.Post_Id,
		User_Id: request.User_Id,
		Content: request.Content,
	}

	createdComment := service.repository.Save(ctx, tx, newComment)

	return web.ToCommentResponse(createdComment)
}

func (service *CommentServiceImpl) FindCommentsByPost(ctx context.Context, postId, limit, offset int) ([]web.CommentResponse, int) {
	tx, err := service.db.Begin()
	helpers.PanicError(err)
	defer helpers.TxRollbackCommit(tx)

	comments := service.repository.FindCommentsByPost(ctx, tx, postId, limit, offset)
	countComments := service.repository.CountCommentsByPost(ctx, tx, postId)

	var commentsData []web.CommentResponse

	if len(comments) == 0 {
		return commentsData, 0
	}

	for _, comment := range comments {
		commentsData = append(commentsData, web.ToCommentResponse(comment))
	}

	return commentsData, countComments
}

func (service *CommentServiceImpl) FindById(ctx context.Context, commentId int) web.CommentResponse {
	tx, err := service.db.Begin()
	helpers.PanicError(err)
	defer helpers.TxRollbackCommit(tx)

	comment := service.repository.FindById(ctx, tx, commentId)

	return web.ToCommentResponse(comment)
}

func (service *CommentServiceImpl) Update(ctx context.Context, request web.CommentUpdateRequest) web.CommentResponse {
	err := service.validator.Struct(request)
	helpers.PanicError(err)

	tx, err := service.db.Begin()
	helpers.PanicError(err)
	defer helpers.TxRollbackCommit(tx)

	service.repository.FindById(ctx, tx, request.Id)

	updatedCommentData := domain.Comment{
		Id:      request.Id,
		Content: request.Content,
	}

	updatedComment := service.repository.Update(ctx, tx, updatedCommentData)

	return web.ToCommentResponse(updatedComment)
}

func (service *CommentServiceImpl) Delete(ctx context.Context, commentId int) {
	tx, err := service.db.Begin()
	helpers.PanicError(err)
	defer helpers.TxRollbackCommit(tx)

	service.repository.Delete(ctx, tx, commentId)
}
