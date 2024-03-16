package controllers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type UserController interface {
	CreateUser(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	SignInUser(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindAllUser(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindByIdUser(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	UpdateUser(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	DeleteUser(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	GetRefreshToken(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
