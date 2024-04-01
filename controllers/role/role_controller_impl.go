package role

import (
	"net/http"
	"strconv"

	"github.com/hutamatr/GoBlogify/helpers"
	"github.com/hutamatr/GoBlogify/model/web"
	servicesRole "github.com/hutamatr/GoBlogify/services/role"
	"github.com/julienschmidt/httprouter"
)

type RoleControllerImpl struct {
	service servicesRole.RoleService
}

func NewRoleController(roleService servicesRole.RoleService) RoleController {
	return &RoleControllerImpl{
		service: roleService,
	}
}

func (controller *RoleControllerImpl) CreateRole(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	var roleRequest web.RoleCreateRequest
	helpers.DecodeJSONFromRequest(request, &roleRequest)

	isAdmin := helpers.IsAdmin(request)

	role := controller.service.Create(request.Context(), roleRequest, isAdmin)

	roleResponse := web.ResponseJSON{
		Code:   http.StatusCreated,
		Status: "CREATED",
		Data:   role,
	}
	writer.WriteHeader(http.StatusCreated)

	helpers.EncodeJSONFromResponse(writer, roleResponse)
}

func (controller *RoleControllerImpl) FindAllRole(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	isAdmin := helpers.IsAdmin(request)

	roles := controller.service.FindAll(request.Context(), isAdmin)

	roleResponse := web.ResponseJSON{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   roles,
	}

	helpers.EncodeJSONFromResponse(writer, roleResponse)
}

func (controller *RoleControllerImpl) FindRoleById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	isAdmin := helpers.IsAdmin(request)

	id := params.ByName("roleId")
	roleId, err := strconv.Atoi(id)

	helpers.PanicError(err)

	role := controller.service.FindById(request.Context(), roleId, isAdmin)

	roleResponse := web.ResponseJSON{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   role,
	}

	helpers.EncodeJSONFromResponse(writer, roleResponse)
}

func (controller *RoleControllerImpl) UpdateRole(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	var roleUpdateRequest web.RoleUpdateRequest
	helpers.DecodeJSONFromRequest(request, &roleUpdateRequest)

	isAdmin := helpers.IsAdmin(request)

	id := params.ByName("roleId")
	roleId, err := strconv.Atoi(id)
	helpers.PanicError(err)

	roleUpdateRequest.Id = roleId

	updatedRole := controller.service.Update(request.Context(), roleUpdateRequest, isAdmin)

	roleResponse := web.ResponseJSON{
		Code:   http.StatusOK,
		Status: "UPDATED",
		Data:   updatedRole,
	}

	helpers.EncodeJSONFromResponse(writer, roleResponse)
}

func (controller *RoleControllerImpl) DeleteRole(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	isAdmin := helpers.IsAdmin(request)

	id := params.ByName("roleId")
	roleId, err := strconv.Atoi(id)
	helpers.PanicError(err)

	controller.service.Delete(request.Context(), roleId, isAdmin)

	roleResponse := web.ResponseJSON{
		Code:   http.StatusOK,
		Status: "DELETED",
	}

	helpers.EncodeJSONFromResponse(writer, roleResponse)
}
