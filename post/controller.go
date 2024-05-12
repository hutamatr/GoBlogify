package post

import (
	"net/http"
	"strconv"

	"github.com/hutamatr/GoBlogify/helpers"
	"github.com/julienschmidt/httprouter"
)

type PostController interface {
	CreatePostHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindAllPostByUserHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindAllPostByFollowedHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindByIdPostHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	UpdatePostHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	DeletePostHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}

type PostControllerImpl struct {
	service PostService
}

func NewPostController(postService PostService) PostController {
	return &PostControllerImpl{
		service: postService,
	}
}

func (controller *PostControllerImpl) CreatePostHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	var postRequest PostCreateRequest
	helpers.DecodeJSONFromRequest(request, &postRequest)

	post := controller.service.Create(request.Context(), postRequest)

	postResponse := helpers.ResponseJSON{
		Code:   http.StatusCreated,
		Status: "CREATED",
		Data:   post,
	}

	writer.WriteHeader(http.StatusCreated)
	helpers.EncodeJSONFromResponse(writer, postResponse)
}

func (controller *PostControllerImpl) FindAllPostByUserHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id := params.ByName("userId")
	userId, err := strconv.Atoi(id)
	helpers.PanicError(err, "Invalid User Id")
	limit, offset := helpers.GetLimitOffset(request)

	posts, countPosts := controller.service.FindAllByUser(request.Context(), userId, limit, offset)

	postResponse := helpers.ResponseJSON{
		Code:   http.StatusOK,
		Status: "OK",
		Data: map[string]interface{}{
			"posts":  posts,
			"limit":  limit,
			"offset": offset,
			"total":  countPosts,
		},
	}

	writer.WriteHeader(http.StatusOK)
	helpers.EncodeJSONFromResponse(writer, postResponse)
}

func (controller *PostControllerImpl) FindAllPostByFollowedHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id := params.ByName("userId")
	userId, err := strconv.Atoi(id)
	helpers.PanicError(err, "Invalid User Id")
	limit, offset := helpers.GetLimitOffset(request)

	postsByFollowed, countPosts := controller.service.FindAllByFollowed(request.Context(), userId, limit, offset)

	postResponse := helpers.ResponseJSON{
		Code:   http.StatusOK,
		Status: "OK",
		Data: map[string]interface{}{
			"posts":  postsByFollowed,
			"limit":  limit,
			"offset": offset,
			"total":  countPosts,
		},
	}

	writer.WriteHeader(http.StatusOK)
	helpers.EncodeJSONFromResponse(writer, postResponse)
}

func (controller *PostControllerImpl) FindByIdPostHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id := params.ByName("postId")
	postId, err := strconv.Atoi(id)

	helpers.PanicError(err, "Invalid Post Id")

	post := controller.service.FindById(request.Context(), postId)

	postResponse := helpers.ResponseJSON{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   post,
	}

	writer.WriteHeader(http.StatusOK)
	helpers.EncodeJSONFromResponse(writer, postResponse)
}

func (controller *PostControllerImpl) UpdatePostHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id := params.ByName("postId")
	postId, err := strconv.Atoi(id)

	helpers.PanicError(err, "Invalid Post Id")

	var postUpdateRequest PostUpdateRequest

	postUpdateRequest.Id = postId

	helpers.DecodeJSONFromRequest(request, &postUpdateRequest)

	updatedPost := controller.service.Update(request.Context(), postUpdateRequest)

	postResponse := helpers.ResponseJSON{
		Code:   http.StatusOK,
		Status: "UPDATED",
		Data:   updatedPost,
	}

	writer.WriteHeader(http.StatusOK)
	helpers.EncodeJSONFromResponse(writer, postResponse)
}

func (controller *PostControllerImpl) DeletePostHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id := params.ByName("postId")
	postId, err := strconv.Atoi(id)

	helpers.PanicError(err, "Invalid Post Id")

	controller.service.Delete(request.Context(), postId)

	postResponse := helpers.ResponseJSON{
		Code:   http.StatusOK,
		Status: "DELETED",
	}

	writer.WriteHeader(http.StatusOK)
	helpers.EncodeJSONFromResponse(writer, postResponse)
}
