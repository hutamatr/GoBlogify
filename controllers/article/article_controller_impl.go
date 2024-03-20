package controllers

import (
	"net/http"
	"strconv"

	"github.com/hutamatr/go-blog-api/helpers"
	"github.com/hutamatr/go-blog-api/model/web"
	servicesA "github.com/hutamatr/go-blog-api/services/article"
	"github.com/julienschmidt/httprouter"
)

type ArticleControllerImpl struct {
	service servicesA.ArticleService
}

func NewArticleController(articleService servicesA.ArticleService) ArticleController {
	return &ArticleControllerImpl{
		service: articleService,
	}
}

func (controller *ArticleControllerImpl) CreateArticle(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	var articleRequest web.ArticleCreateRequest
	helpers.DecodeJSONFromRequest(request, &articleRequest)

	article := controller.service.Create(request.Context(), articleRequest)

	writer.WriteHeader(http.StatusCreated)
	articleResponse := web.ResponseJSON{
		Code:   http.StatusCreated,
		Status: "CREATED",
		Data:   article,
	}

	helpers.EncodeJSONFromResponse(writer, articleResponse)
}

func (controller *ArticleControllerImpl) FindAllArticle(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	articles := controller.service.FindAll(request.Context())

	articleResponse := web.ResponseJSON{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   articles,
	}

	helpers.EncodeJSONFromResponse(writer, articleResponse)
}

func (controller *ArticleControllerImpl) FindByIdArticle(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id := params.ByName("articleId")
	articleId, err := strconv.Atoi(id)

	helpers.PanicError(err)

	article := controller.service.FindById(request.Context(), articleId)

	articleResponse := web.ResponseJSON{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   article,
	}

	helpers.EncodeJSONFromResponse(writer, articleResponse)
}

func (controller *ArticleControllerImpl) UpdateArticle(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	id := params.ByName("articleId")
	articleId, err := strconv.Atoi(id)

	helpers.PanicError(err)

	var articleUpdateRequest web.ArticleUpdateRequest

	articleUpdateRequest.Id = articleId

	helpers.DecodeJSONFromRequest(request, &articleUpdateRequest)

	updatedArticle := controller.service.Update(request.Context(), articleUpdateRequest)

	articleResponse := web.ResponseJSON{
		Code:   http.StatusOK,
		Status: "UPDATED",
		Data:   updatedArticle,
	}

	helpers.EncodeJSONFromResponse(writer, articleResponse)
}

func (controller *ArticleControllerImpl) DeleteArticle(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id := params.ByName("articleId")
	articleId, err := strconv.Atoi(id)

	helpers.PanicError(err)

	controller.service.Delete(request.Context(), articleId)

	articleResponse := web.ResponseJSON{
		Code:   http.StatusOK,
		Status: "DELETED",
	}

	helpers.EncodeJSONFromResponse(writer, articleResponse)
}
