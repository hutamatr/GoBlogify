package controllers

import (
	"net/http"
	"strconv"

	"github.com/hutamatr/GoBlogify/helpers"
	"github.com/hutamatr/GoBlogify/model/web"
	servicesComment "github.com/hutamatr/GoBlogify/services/comment"
	"github.com/julienschmidt/httprouter"
)

type CommentControllerImpl struct {
	service servicesComment.CommentService
}

func NewCommentController(service servicesComment.CommentService) CommentController {
	return &CommentControllerImpl{
		service: service,
	}
}

func (controller *CommentControllerImpl) CreateComment(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	var CommentRequest web.CommentCreateRequest
	helpers.DecodeJSONFromRequest(request, &CommentRequest)

	comment := controller.service.Create(request.Context(), CommentRequest)

	CommentResponse := web.ResponseJSON{
		Code:   http.StatusCreated,
		Status: "CREATED",
		Data:   comment,
	}

	writer.WriteHeader(http.StatusCreated)

	helpers.EncodeJSONFromResponse(writer, CommentResponse)
}

func (controller *CommentControllerImpl) FindCommentsByPost(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id := request.URL.Query().Get("postId")
	postId, err := strconv.Atoi(id)
	helpers.PanicError(err)

	offset := request.URL.Query().Get("offset")
	if offset == "" {
		offset = "0"
	}
	offsetVal, err := strconv.Atoi(offset)
	helpers.PanicError(err)

	comments := controller.service.FindCommentsByPost(request.Context(), postId, offsetVal)

	CommentResponse := web.ResponseJSON{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   comments,
	}

	helpers.EncodeJSONFromResponse(writer, CommentResponse)
}

func (controller *CommentControllerImpl) FindByIdComment(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id := params.ByName("commentId")
	commentId, err := strconv.Atoi(id)

	helpers.PanicError(err)

	comment := controller.service.FindById(request.Context(), commentId)

	CommentResponse := web.ResponseJSON{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   comment,
	}

	helpers.EncodeJSONFromResponse(writer, CommentResponse)
}

func (controller *CommentControllerImpl) UpdateComment(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	var CommentUpdateRequest web.CommentUpdateRequest
	helpers.DecodeJSONFromRequest(request, &CommentUpdateRequest)

	id := params.ByName("commentId")
	commentId, err := strconv.Atoi(id)
	helpers.PanicError(err)

	CommentUpdateRequest.Id = commentId

	updatedComment := controller.service.Update(request.Context(), CommentUpdateRequest)

	CommentResponse := web.ResponseJSON{
		Code:   http.StatusOK,
		Status: "UPDATED",
		Data:   updatedComment,
	}

	helpers.EncodeJSONFromResponse(writer, CommentResponse)
}

func (controller *CommentControllerImpl) DeleteComment(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id := params.ByName("commentId")
	commentId, err := strconv.Atoi(id)

	helpers.PanicError(err)

	controller.service.Delete(request.Context(), commentId)

	CommentResponse := web.ResponseJSON{
		Code:   http.StatusOK,
		Status: "DELETED",
	}

	helpers.EncodeJSONFromResponse(writer, CommentResponse)
}
