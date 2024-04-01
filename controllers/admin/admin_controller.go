package controllers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type AdminController interface {
	CreateAdmin(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	SignInAdmin(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
