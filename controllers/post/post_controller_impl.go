package controllers

import (
	"net/http"
	"strconv"

	"github.com/hutamatr/GoBlogify/helpers"
	"github.com/hutamatr/GoBlogify/model/web"
	servicesPost "github.com/hutamatr/GoBlogify/services/post"
	"github.com/julienschmidt/httprouter"
)

type PostControllerImpl struct {
	service servicesPost.PostService
}

func NewPostController(postService servicesPost.PostService) PostController {
	return &PostControllerImpl{
		service: postService,
	}
}

func (controller *PostControllerImpl) CreatePost(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	var postRequest web.PostCreateRequest
	helpers.DecodeJSONFromRequest(request, &postRequest)

	post := controller.service.Create(request.Context(), postRequest)

	postResponse := web.ResponseJSON{
		Code:   http.StatusCreated,
		Status: "CREATED",
		Data:   post,
	}
	writer.WriteHeader(http.StatusCreated)

	helpers.EncodeJSONFromResponse(writer, postResponse)
}

func (controller *PostControllerImpl) FindAllPost(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	posts := controller.service.FindAll(request.Context())

	postResponse := web.ResponseJSON{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   posts,
	}

	helpers.EncodeJSONFromResponse(writer, postResponse)
}

func (controller *PostControllerImpl) FindByIdPost(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id := params.ByName("postId")
	postId, err := strconv.Atoi(id)

	helpers.PanicError(err)

	post := controller.service.FindById(request.Context(), postId)

	postResponse := web.ResponseJSON{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   post,
	}

	helpers.EncodeJSONFromResponse(writer, postResponse)
}

func (controller *PostControllerImpl) UpdatePost(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	id := params.ByName("postId")
	postId, err := strconv.Atoi(id)

	helpers.PanicError(err)

	var postUpdateRequest web.PostUpdateRequest

	postUpdateRequest.Id = postId

	helpers.DecodeJSONFromRequest(request, &postUpdateRequest)

	updatedPost := controller.service.Update(request.Context(), postUpdateRequest)

	postResponse := web.ResponseJSON{
		Code:   http.StatusOK,
		Status: "UPDATED",
		Data:   updatedPost,
	}

	helpers.EncodeJSONFromResponse(writer, postResponse)
}

func (controller *PostControllerImpl) DeletePost(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id := params.ByName("postId")
	postId, err := strconv.Atoi(id)

	helpers.PanicError(err)

	controller.service.Delete(request.Context(), postId)

	postResponse := web.ResponseJSON{
		Code:   http.StatusOK,
		Status: "DELETED",
	}

	helpers.EncodeJSONFromResponse(writer, postResponse)
}
