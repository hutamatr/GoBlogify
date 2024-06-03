package post

import (
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/hutamatr/GoBlogify/exception"
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
	err := request.ParseMultipartForm(10 << 20) // 10 MB

	if err != nil {
		panic(exception.NewBadRequestError(err.Error()))
	}

	multiPartFormData := request.MultipartForm

	title := multiPartFormData.Value["title"][0]
	postBody := multiPartFormData.Value["post_body"][0]
	published := multiPartFormData.Value["published"][0]
	userId := multiPartFormData.Value["user_id"][0]
	categoryId := multiPartFormData.Value["category_id"][0]

	attachments := []map[string]interface{}{}

	for _, v := range multiPartFormData.File["attachment"] {
		file, err := v.Open()
		helpers.PanicError(err, "failed open file")

		attachment := map[string]interface{}{
			"image":    file,
			"filename": v.Filename,
		}
		defer file.Close()

		attachments = append(attachments, attachment)
	}

	userID, err := strconv.Atoi(userId)
	helpers.PanicError(err, "Failed convert user_id")
	categoryID, err := strconv.Atoi(categoryId)
	helpers.PanicError(err, "Failed convert category_id")

	var isPublished bool
	switch published {
	case "true":
		isPublished = true
	default:
		isPublished = false
	}

	var postRequest = PostCreateRequest{
		Title:        title,
		Post_Body:    postBody,
		Published:    isPublished,
		User_Id:      userID,
		Category_Id:  categoryID,
		Image_1:      attachments[0]["image"].(multipart.File),
		Image_Name_1: attachments[0]["filename"].(string),
		Image_2:      attachments[1]["image"].(multipart.File),
		Image_Name_2: attachments[1]["filename"].(string),
		Image_3:      attachments[2]["image"].(multipart.File),
		Image_Name_3: attachments[2]["filename"].(string),
	}

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

	request.ParseMultipartForm(10 << 20) // 10 MB

	multiPartFormData := request.MultipartForm

	title := multiPartFormData.Value["title"][0]
	postBody := multiPartFormData.Value["post_body"][0]
	published := multiPartFormData.Value["published"][0]
	userId := multiPartFormData.Value["user_id"][0]
	categoryId := multiPartFormData.Value["category_id"][0]

	attachments := map[int]map[string]interface{}{}

	for i, v := range multiPartFormData.File["attachment"] {
		file, err := v.Open()
		helpers.PanicError(err, "failed open file")

		attachment := map[string]interface{}{
			"image":    file,
			"filename": v.Filename,
		}

		attachments[i] = attachment
	}

	userID, err := strconv.Atoi(userId)
	helpers.PanicError(err, "Failed convert user_id")
	categoryID, err := strconv.Atoi(categoryId)
	helpers.PanicError(err, "Failed convert category_id")

	var isPublished bool
	switch published {
	case "true":
		isPublished = true
	default:
		isPublished = false
	}

	var postUpdateRequest PostUpdateRequest

	postUpdateRequest.Id = postId
	postUpdateRequest.Title = title
	postUpdateRequest.Post_Body = postBody
	postUpdateRequest.Published = isPublished
	postUpdateRequest.User_Id = userID
	postUpdateRequest.Category_Id = categoryID

	if attachment, exists := attachments[0]; exists && attachment["filename"].(string) != "" {
		postUpdateRequest.Image_1 = attachment["image"].(multipart.File)
		postUpdateRequest.Image_Name_1 = attachment["filename"].(string)
	}
	if attachment, exists := attachments[1]; exists && attachment["filename"].(string) != "" {
		postUpdateRequest.Image_2 = attachment["image"].(multipart.File)
		postUpdateRequest.Image_Name_2 = attachment["filename"].(string)
	}
	if attachment, exists := attachments[2]; exists && attachment["filename"].(string) != "" {
		postUpdateRequest.Image_3 = attachment["image"].(multipart.File)
		postUpdateRequest.Image_Name_3 = attachment["filename"].(string)
	}

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
