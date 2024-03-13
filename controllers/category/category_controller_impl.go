package controllers

import (
	"net/http"
	"strconv"

	"github.com/hutamatr/go-blog-api/helpers"
	"github.com/hutamatr/go-blog-api/model/web"
	servicesC "github.com/hutamatr/go-blog-api/services/category"
	"github.com/julienschmidt/httprouter"
)

type CategoryControllerImpl struct {
	service servicesC.CategoryService
}

func NewCategoryController(categoryService servicesC.CategoryService) CategoryController {
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
		Status: "OK",
		Data:   category,
	}

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
	id := params.ByName("CategoryId")
	CategoryId, err := strconv.Atoi(id)

	helpers.PanicError(err)

	category := controller.service.FindById(request.Context(), CategoryId)

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

	updatedCategory := controller.service.Update(request.Context(), CategoryUpdateRequest)

	CategoryResponse := web.ResponseJSON{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   updatedCategory,
	}

	helpers.EncodeJSONFromResponse(writer, CategoryResponse)
}

func (controller *CategoryControllerImpl) DeleteCategory(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id := params.ByName("CategoryId")
	CategoryId, err := strconv.Atoi(id)

	helpers.PanicError(err)

	controller.service.Delete(request.Context(), CategoryId)

	CategoryResponse := web.ResponseJSON{
		Code:   http.StatusOK,
		Status: "OK",
	}

	helpers.EncodeJSONFromResponse(writer, CategoryResponse)
}
