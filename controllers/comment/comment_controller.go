package controllers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type CommentController interface {
	CreateComment(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindCommentsByPost(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindByIdComment(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	UpdateComment(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	DeleteComment(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
