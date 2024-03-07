package go_blog

import (
	"net/http"
	"strconv"

	"github.com/hutamatr/go-blog-api/cmd/go_blog/helper"
	"github.com/julienschmidt/httprouter"
)

type CategoryController interface {
	CreateCategory(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindAllCategory(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindByIdCategory(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	UpdateCategory(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	DeleteCategory(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}

type CategoryControllerImpl struct {
	service CategoryService
}

func (controller *CategoryControllerImpl) CreateCategory(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	CategoryRequest := CategoryCreateRequest{}
	helper.DecodeJSONFromRequest(request, &CategoryRequest)

	category := controller.service.Create(request.Context(), CategoryRequest)

	CategoryResponse := ResponseJSON{
		Code:   http.StatusCreated,
		Status: "OK",
		Data:   category,
	}

	helper.EncodeJSONFromResponse(writer, CategoryResponse)
}

func (controller *CategoryControllerImpl) FindAllCategory(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	categories := controller.service.FindAll(request.Context())

	CategoryResponse := ResponseJSON{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   categories,
	}

	helper.EncodeJSONFromResponse(writer, CategoryResponse)
}

func (controller *CategoryControllerImpl) FindByIdCategory(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id := params.ByName("CategoryId")
	CategoryId, err := strconv.Atoi(id)

	helper.PanicError(err)

	category := controller.service.FindById(request.Context(), CategoryId)

	CategoryResponse := ResponseJSON{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   category,
	}

	helper.EncodeJSONFromResponse(writer, CategoryResponse)
}

func (controller *CategoryControllerImpl) UpdateCategory(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	CategoryUpdateRequest := CategoryUpdateRequest{}

	helper.DecodeJSONFromRequest(request, &CategoryUpdateRequest)

	updatedCategory := controller.service.Update(request.Context(), CategoryUpdateRequest)

	CategoryResponse := ResponseJSON{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   updatedCategory,
	}

	helper.EncodeJSONFromResponse(writer, CategoryResponse)
}

func (controller *CategoryControllerImpl) DeleteCategory(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	id := params.ByName("CategoryId")

	CategoryId, err := strconv.Atoi(id)
	helper.PanicError(err)

	controller.service.Delete(request.Context(), CategoryId)

	CategoryResponse := ResponseJSON{
		Code:   http.StatusOK,
		Status: "OK",
	}

	helper.EncodeJSONFromResponse(writer, CategoryResponse)
}
