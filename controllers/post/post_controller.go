package controllers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type PostController interface {
	CreatePost(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindAllPost(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindByIdPost(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	UpdatePost(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	DeletePost(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
