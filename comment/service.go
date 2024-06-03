package comment

import (
	"context"
	"database/sql"

	"github.com/go-playground/validator/v10"
	"github.com/hutamatr/GoBlogify/exception"
	"github.com/hutamatr/GoBlogify/helpers"
)

type CommentService interface {
	Create(ctx context.Context, request CommentCreateRequest) CommentResponse
	FindCommentsByPost(ctx context.Context, postId, limit, offset int) ([]CommentResponse, int)
	FindById(ctx context.Context, commentId int) CommentResponse
	Update(ctx context.Context, request CommentUpdateRequest) CommentResponse
	Delete(ctx context.Context, commentId int)
}

type CommentServiceImpl struct {
	repository CommentRepository
	db         *sql.DB
	validator  *validator.Validate
}

func NewCommentService(commentRepository CommentRepository, db *sql.DB, validator *validator.Validate) CommentService {
	return &CommentServiceImpl{
		repository: commentRepository,
		db:         db,
		validator:  validator,
	}
}

func (service *CommentServiceImpl) Create(ctx context.Context, request CommentCreateRequest) CommentResponse {
	err := service.validator.Struct(request)
	helpers.PanicError(err, "invalid request")

	tx, err := service.db.Begin()
	helpers.PanicError(err, "failed to begin transaction")
	defer helpers.TxRollbackCommit(tx)

	newComment := Comment{
		Post_Id: request.Post_Id,
		User_Id: request.User_Id,
		Comment: request.Comment,
	}

	createdComment := service.repository.Save(ctx, tx, newComment)

	return ToCommentResponse(createdComment)
}

func (service *CommentServiceImpl) FindCommentsByPost(ctx context.Context, postId, limit, offset int) ([]CommentResponse, int) {
	tx, err := service.db.Begin()
	helpers.PanicError(err, "failed to begin transaction")
	defer helpers.TxRollbackCommit(tx)

	comments := service.repository.FindCommentsByPost(ctx, tx, postId, limit, offset)
	countComments := service.repository.CountCommentsByPost(ctx, tx, postId)

	var commentsData []CommentResponse

	if len(comments) == 0 {
		panic(exception.NewNotFoundError("comments not found"))
	}

	for _, comment := range comments {
		commentsData = append(commentsData, ToCommentResponse(comment))
	}

	return commentsData, countComments
}

func (service *CommentServiceImpl) FindById(ctx context.Context, commentId int) CommentResponse {
	tx, err := service.db.Begin()
	helpers.PanicError(err, "failed to begin transaction")
	defer helpers.TxRollbackCommit(tx)

	comment := service.repository.FindById(ctx, tx, commentId)

	return ToCommentResponse(comment)
}

func (service *CommentServiceImpl) Update(ctx context.Context, request CommentUpdateRequest) CommentResponse {
	err := service.validator.Struct(request)
	helpers.PanicError(err, "invalid request")

	tx, err := service.db.Begin()
	helpers.PanicError(err, "failed to begin transaction")
	defer helpers.TxRollbackCommit(tx)

	service.repository.FindById(ctx, tx, request.Id)

	updatedCommentData := Comment{
		Id:      request.Id,
		Comment: request.Comment,
	}

	updatedComment := service.repository.Update(ctx, tx, updatedCommentData)

	return ToCommentResponse(updatedComment)
}

func (service *CommentServiceImpl) Delete(ctx context.Context, commentId int) {
	tx, err := service.db.Begin()
	helpers.PanicError(err, "failed to begin transaction")
	defer helpers.TxRollbackCommit(tx)

	service.repository.Delete(ctx, tx, commentId)
}
