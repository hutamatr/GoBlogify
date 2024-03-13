package controllers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type CategoryController interface {
	CreateCategory(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindAllCategory(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindByIdCategory(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	UpdateCategory(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	DeleteCategory(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
