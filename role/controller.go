package role

import (
	"net/http"
	"strconv"

	"github.com/hutamatr/GoBlogify/helpers"
	"github.com/julienschmidt/httprouter"
)

type RoleController interface {
	CreateRoleHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindAllRoleHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindRoleByIdHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	UpdateRoleHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	DeleteRoleHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}

type RoleControllerImpl struct {
	service RoleService
}

func NewRoleController(roleService RoleService) RoleController {
	return &RoleControllerImpl{
		service: roleService,
	}
}

func (controller *RoleControllerImpl) CreateRoleHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	var roleRequest RoleCreateRequest
	helpers.DecodeJSONFromRequest(request, &roleRequest)

	isAdmin := helpers.IsAdmin(request)

	role := controller.service.Create(request.Context(), roleRequest, isAdmin)

	roleResponse := helpers.ResponseJSON{
		Code:   http.StatusCreated,
		Status: "CREATED",
		Data:   role,
	}

	writer.WriteHeader(http.StatusCreated)
	helpers.EncodeJSONFromResponse(writer, roleResponse)
}

func (controller *RoleControllerImpl) FindAllRoleHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	isAdmin := helpers.IsAdmin(request)

	roles := controller.service.FindAll(request.Context(), isAdmin)

	roleResponse := helpers.ResponseJSON{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   roles,
	}

	writer.WriteHeader(http.StatusOK)
	helpers.EncodeJSONFromResponse(writer, roleResponse)
}

func (controller *RoleControllerImpl) FindRoleByIdHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	isAdmin := helpers.IsAdmin(request)

	id := params.ByName("roleId")
	roleId, err := strconv.Atoi(id)

	helpers.PanicError(err, "Invalid Role Id")

	role := controller.service.FindById(request.Context(), roleId, isAdmin)

	roleResponse := helpers.ResponseJSON{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   role,
	}

	writer.WriteHeader(http.StatusOK)
	helpers.EncodeJSONFromResponse(writer, roleResponse)
}

func (controller *RoleControllerImpl) UpdateRoleHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	var roleUpdateRequest RoleUpdateRequest
	helpers.DecodeJSONFromRequest(request, &roleUpdateRequest)

	isAdmin := helpers.IsAdmin(request)

	id := params.ByName("roleId")
	roleId, err := strconv.Atoi(id)
	helpers.PanicError(err, "Invalid Role Id")

	roleUpdateRequest.Id = roleId

	updatedRole := controller.service.Update(request.Context(), roleUpdateRequest, isAdmin)

	roleResponse := helpers.ResponseJSON{
		Code:   http.StatusOK,
		Status: "UPDATED",
		Data:   updatedRole,
	}

	writer.WriteHeader(http.StatusOK)
	helpers.EncodeJSONFromResponse(writer, roleResponse)
}

func (controller *RoleControllerImpl) DeleteRoleHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	isAdmin := helpers.IsAdmin(request)

	id := params.ByName("roleId")
	roleId, err := strconv.Atoi(id)
	helpers.PanicError(err, "Invalid Role Id")

	controller.service.Delete(request.Context(), roleId, isAdmin)

	roleResponse := helpers.ResponseJSON{
		Code:   http.StatusOK,
		Status: "DELETED",
	}

	writer.WriteHeader(http.StatusOK)
	helpers.EncodeJSONFromResponse(writer, roleResponse)
}
