package role

import (
	"net/http"
	"strconv"

	"github.com/hutamatr/go-blog-api/helpers"
	"github.com/hutamatr/go-blog-api/model/web"
	servicesR "github.com/hutamatr/go-blog-api/services/role"
	"github.com/julienschmidt/httprouter"
)

type RoleControllerImpl struct {
	service servicesR.RoleService
}

func NewRoleController(roleService servicesR.RoleService) RoleController {
	return &RoleControllerImpl{
		service: roleService,
	}
}

func (controller *RoleControllerImpl) CreateRole(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	var roleRequest web.RoleCreateRequest
	helpers.DecodeJSONFromRequest(request, &roleRequest)

	role := controller.service.Create(request.Context(), roleRequest)

	roleResponse := web.ResponseJSON{
		Code:   http.StatusCreated,
		Status: "CREATED",
		Data:   role,
	}

	helpers.EncodeJSONFromResponse(writer, roleResponse)
}

func (controller *RoleControllerImpl) FindAllRole(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	roles := controller.service.FindAll(request.Context())

	roleResponse := web.ResponseJSON{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   roles,
	}

	helpers.EncodeJSONFromResponse(writer, roleResponse)
}

func (controller *RoleControllerImpl) FindRoleById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id := params.ByName("roleId")
	roleId, err := strconv.Atoi(id)

	helpers.PanicError(err)

	role := controller.service.FindById(request.Context(), roleId)

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

	id := params.ByName("roleId")
	roleId, err := strconv.Atoi(id)
	helpers.PanicError(err)

	roleUpdateRequest.Id = roleId

	updatedRole := controller.service.Update(request.Context(), roleUpdateRequest)

	roleResponse := web.ResponseJSON{
		Code:   http.StatusOK,
		Status: "UPDATED",
		Data:   updatedRole,
	}

	helpers.EncodeJSONFromResponse(writer, roleResponse)
}

func (controller *RoleControllerImpl) DeleteRole(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id := params.ByName("roleId")
	roleId, err := strconv.Atoi(id)
	helpers.PanicError(err)

	controller.service.Delete(request.Context(), roleId)

	roleResponse := web.ResponseJSON{
		Code:   http.StatusOK,
		Status: "DELETED",
	}

	helpers.EncodeJSONFromResponse(writer, roleResponse)
}
