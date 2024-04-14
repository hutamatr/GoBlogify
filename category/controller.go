package category

import (
	"net/http"
	"strconv"

	"github.com/hutamatr/GoBlogify/helpers"
	"github.com/julienschmidt/httprouter"
)

type CategoryController interface {
	CreateCategoryHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindAllCategoryHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindByIdCategoryHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	UpdateCategoryHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	DeleteCategoryHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}

type CategoryControllerImpl struct {
	service CategoryService
}

func NewCategoryController(categoryService CategoryService) CategoryController {
	return &CategoryControllerImpl{
		service: categoryService,
	}
}

func (controller *CategoryControllerImpl) CreateCategoryHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	var CategoryRequest CategoryCreateRequest
	helpers.DecodeJSONFromRequest(request, &CategoryRequest)

	isAdmin := helpers.IsAdmin(request)

	category := controller.service.Create(request.Context(), CategoryRequest, isAdmin)

	CategoryResponse := helpers.ResponseJSON{
		Code:   http.StatusCreated,
		Status: "CREATED",
		Data:   category,
	}

	writer.WriteHeader(http.StatusCreated)
	helpers.EncodeJSONFromResponse(writer, CategoryResponse)
}

func (controller *CategoryControllerImpl) FindAllCategoryHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	limit, offset := helpers.GetLimitOffset(request)

	categories, countCategories := controller.service.FindAll(request.Context(), limit, offset)

	CategoryResponse := helpers.ResponseJSON{
		Code:   http.StatusOK,
		Status: "OK",
		Data: map[string]interface{}{
			"categories": categories,
			"limit":      limit,
			"offset":     offset,
			"total":      countCategories,
		},
	}

	writer.WriteHeader(http.StatusOK)
	helpers.EncodeJSONFromResponse(writer, CategoryResponse)
}

func (controller *CategoryControllerImpl) FindByIdCategoryHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id := params.ByName("categoryId")
	categoryId, err := strconv.Atoi(id)

	helpers.PanicError(err, "Invalid Category Id")

	category := controller.service.FindById(request.Context(), categoryId)

	CategoryResponse := helpers.ResponseJSON{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   category,
	}

	writer.WriteHeader(http.StatusOK)
	helpers.EncodeJSONFromResponse(writer, CategoryResponse)
}

func (controller *CategoryControllerImpl) UpdateCategoryHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	var CategoryUpdateRequest CategoryUpdateRequest
	helpers.DecodeJSONFromRequest(request, &CategoryUpdateRequest)

	isAdmin := helpers.IsAdmin(request)

	id := params.ByName("categoryId")
	categoryId, err := strconv.Atoi(id)
	helpers.PanicError(err, "Invalid Category Id")

	CategoryUpdateRequest.Id = categoryId

	updatedCategory := controller.service.Update(request.Context(), CategoryUpdateRequest, isAdmin)

	CategoryResponse := helpers.ResponseJSON{
		Code:   http.StatusOK,
		Status: "UPDATED",
		Data:   updatedCategory,
	}

	writer.WriteHeader(http.StatusOK)
	helpers.EncodeJSONFromResponse(writer, CategoryResponse)
}

func (controller *CategoryControllerImpl) DeleteCategoryHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	isAdmin := helpers.IsAdmin(request)

	id := params.ByName("categoryId")
	categoryId, err := strconv.Atoi(id)

	helpers.PanicError(err, "Invalid Category Id")

	controller.service.Delete(request.Context(), categoryId, isAdmin)

	CategoryResponse := helpers.ResponseJSON{
		Code:   http.StatusOK,
		Status: "DELETED",
	}

	writer.WriteHeader(http.StatusOK)
	helpers.EncodeJSONFromResponse(writer, CategoryResponse)
}
