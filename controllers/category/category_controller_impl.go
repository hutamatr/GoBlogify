package controllers

import (
	"net/http"
	"strconv"

	"github.com/hutamatr/GoBlogify/helpers"
	"github.com/hutamatr/GoBlogify/model/web"
	servicesCategory "github.com/hutamatr/GoBlogify/services/category"
	"github.com/julienschmidt/httprouter"
)

type CategoryControllerImpl struct {
	service servicesCategory.CategoryService
}

func NewCategoryController(categoryService servicesCategory.CategoryService) CategoryController {
	return &CategoryControllerImpl{
		service: categoryService,
	}
}

func (controller *CategoryControllerImpl) CreateCategory(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	var CategoryRequest web.CategoryCreateRequest
	helpers.DecodeJSONFromRequest(request, &CategoryRequest)

	category := controller.service.Create(request.Context(), CategoryRequest)

	CategoryResponse := web.ResponseJSON{
		Code:   http.StatusCreated,
		Status: "CREATED",
		Data:   category,
	}
	writer.WriteHeader(http.StatusCreated)

	helpers.EncodeJSONFromResponse(writer, CategoryResponse)
}

func (controller *CategoryControllerImpl) FindAllCategory(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	categories := controller.service.FindAll(request.Context())

	CategoryResponse := web.ResponseJSON{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   categories,
	}

	helpers.EncodeJSONFromResponse(writer, CategoryResponse)
}

func (controller *CategoryControllerImpl) FindByIdCategory(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id := params.ByName("categoryId")
	categoryId, err := strconv.Atoi(id)

	helpers.PanicError(err)

	category := controller.service.FindById(request.Context(), categoryId)

	CategoryResponse := web.ResponseJSON{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   category,
	}

	helpers.EncodeJSONFromResponse(writer, CategoryResponse)
}

func (controller *CategoryControllerImpl) UpdateCategory(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	var CategoryUpdateRequest web.CategoryUpdateRequest
	helpers.DecodeJSONFromRequest(request, &CategoryUpdateRequest)

	id := params.ByName("categoryId")
	categoryId, err := strconv.Atoi(id)
	helpers.PanicError(err)

	CategoryUpdateRequest.Id = categoryId

	updatedCategory := controller.service.Update(request.Context(), CategoryUpdateRequest)

	CategoryResponse := web.ResponseJSON{
		Code:   http.StatusOK,
		Status: "UPDATED",
		Data:   updatedCategory,
	}

	helpers.EncodeJSONFromResponse(writer, CategoryResponse)
}

func (controller *CategoryControllerImpl) DeleteCategory(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id := params.ByName("categoryId")
	categoryId, err := strconv.Atoi(id)

	helpers.PanicError(err)

	controller.service.Delete(request.Context(), categoryId)

	CategoryResponse := web.ResponseJSON{
		Code:   http.StatusOK,
		Status: "DELETED",
	}

	helpers.EncodeJSONFromResponse(writer, CategoryResponse)
}
