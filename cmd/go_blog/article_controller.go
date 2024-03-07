package go_blog

import (
	"net/http"
	"strconv"

	"github.com/hutamatr/go-blog-api/cmd/go_blog/helper"
	"github.com/julienschmidt/httprouter"
)

type ArticleController interface {
	CreateArticle(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindAllArticle(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindByIdArticle(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	UpdateArticle(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	DeleteArticle(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}

type ArticleControllerImpl struct {
	service ArticleService
}

func (controller *ArticleControllerImpl) CreateArticle(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	articleRequest := ArticleCreateRequest{}
	helper.DecodeJSONFromRequest(request, &articleRequest)

	article := controller.service.Create(request.Context(), articleRequest)

	articleResponse := ResponseJSON{
		Code:   http.StatusCreated,
		Status: "OK",
		Data:   article,
	}

	helper.EncodeJSONFromResponse(writer, articleResponse)
}

func (controller *ArticleControllerImpl) FindAllArticle(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	articles := controller.service.FindAll(request.Context())

	articleResponse := ResponseJSON{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   articles,
	}

	helper.EncodeJSONFromResponse(writer, articleResponse)
}

func (controller *ArticleControllerImpl) FindByIdArticle(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id := params.ByName("articleId")
	articleId, err := strconv.Atoi(id)

	helper.PanicError(err)

	article := controller.service.FindById(request.Context(), articleId)

	articleResponse := ResponseJSON{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   article,
	}

	helper.EncodeJSONFromResponse(writer, articleResponse)
}

func (controller *ArticleControllerImpl) UpdateArticle(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	articleUpdateRequest := ArticleUpdateRequest{}

	helper.DecodeJSONFromRequest(request, &articleUpdateRequest)

	updatedArticle := controller.service.Update(request.Context(), articleUpdateRequest)

	articleResponse := ResponseJSON{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   updatedArticle,
	}

	helper.EncodeJSONFromResponse(writer, articleResponse)
}

func (controller *ArticleControllerImpl) DeleteArticle(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	id := params.ByName("articleId")

	articleId, err := strconv.Atoi(id)
	helper.PanicError(err)

	controller.service.Delete(request.Context(), articleId)

	articleResponse := ResponseJSON{
		Code:   http.StatusOK,
		Status: "OK",
	}

	helper.EncodeJSONFromResponse(writer, articleResponse)
}
