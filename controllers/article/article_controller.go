package controllers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type ArticleController interface {
	CreateArticle(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindAllArticle(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindByIdArticle(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	UpdateArticle(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	DeleteArticle(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
