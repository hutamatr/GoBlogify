package role

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type RoleController interface {
	CreateRole(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindAllRole(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindRoleById(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	UpdateRole(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	DeleteRole(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
