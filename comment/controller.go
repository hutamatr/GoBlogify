package comment

import (
	"net/http"
	"strconv"

	"github.com/hutamatr/GoBlogify/helpers"
	"github.com/julienschmidt/httprouter"
)

type CommentController interface {
	CreateCommentHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindCommentsByPostHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindCommentByIdHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	UpdateCommentHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	DeleteCommentHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}

type CommentControllerImpl struct {
	service CommentService
}

func NewCommentController(service CommentService) CommentController {
	return &CommentControllerImpl{
		service: service,
	}
}

func (controller *CommentControllerImpl) CreateCommentHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	var CommentRequest CommentCreateRequest
	helpers.DecodeJSONFromRequest(request, &CommentRequest)

	comment := controller.service.Create(request.Context(), CommentRequest)

	CommentResponse := helpers.ResponseJSON{
		Code:   http.StatusCreated,
		Status: "CREATED",
		Data:   comment,
	}

	writer.WriteHeader(http.StatusCreated)
	helpers.EncodeJSONFromResponse(writer, CommentResponse)
}

func (controller *CommentControllerImpl) FindCommentsByPostHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id := request.URL.Query().Get("postId")
	postId, err := strconv.Atoi(id)
	helpers.PanicError(err, "Invalid Post Id")

	limit, offset := helpers.GetLimitOffset(request)

	comments, countComments := controller.service.FindCommentsByPost(request.Context(), postId, limit, offset)

	CommentResponse := helpers.ResponseJSON{
		Code:   http.StatusOK,
		Status: "OK",
		Data: map[string]interface{}{
			"comments": comments,
			"limit":    limit,
			"offset":   offset,
			"total":    countComments,
		},
	}

	writer.WriteHeader(http.StatusOK)
	helpers.EncodeJSONFromResponse(writer, CommentResponse)
}

func (controller *CommentControllerImpl) FindCommentByIdHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id := params.ByName("commentId")
	commentId, err := strconv.Atoi(id)

	helpers.PanicError(err, "Invalid Comment Id")

	comment := controller.service.FindById(request.Context(), commentId)

	CommentResponse := helpers.ResponseJSON{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   comment,
	}

	writer.WriteHeader(http.StatusOK)
	helpers.EncodeJSONFromResponse(writer, CommentResponse)
}

func (controller *CommentControllerImpl) UpdateCommentHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	var CommentUpdateRequest CommentUpdateRequest
	helpers.DecodeJSONFromRequest(request, &CommentUpdateRequest)

	id := params.ByName("commentId")
	commentId, err := strconv.Atoi(id)
	helpers.PanicError(err, "Invalid Comment Id")

	CommentUpdateRequest.Id = commentId

	updatedComment := controller.service.Update(request.Context(), CommentUpdateRequest)

	CommentResponse := helpers.ResponseJSON{
		Code:   http.StatusOK,
		Status: "UPDATED",
		Data:   updatedComment,
	}

	writer.WriteHeader(http.StatusOK)
	helpers.EncodeJSONFromResponse(writer, CommentResponse)
}

func (controller *CommentControllerImpl) DeleteCommentHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id := params.ByName("commentId")
	commentId, err := strconv.Atoi(id)

	helpers.PanicError(err, "Invalid Comment Id")

	controller.service.Delete(request.Context(), commentId)

	CommentResponse := helpers.ResponseJSON{
		Code:   http.StatusOK,
		Status: "DELETED",
	}

	writer.WriteHeader(http.StatusOK)
	helpers.EncodeJSONFromResponse(writer, CommentResponse)
}
