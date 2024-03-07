package go_blog

import (
	"net/http"

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
