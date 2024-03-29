package exception

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/hutamatr/go-blog-api/helpers"
	"github.com/hutamatr/go-blog-api/model/web"
)

func ErrorHandler(writer http.ResponseWriter, request *http.Request, err interface{}) {
	if notFoundError(writer, request, err) {
		return
	}
	if validationError(writer, request, err) {
		return
	}
	internalServerError(writer, request, err)
}

func validationError(writer http.ResponseWriter, _ *http.Request, err interface{}) bool {
	if validationError, ok := err.(validator.ValidationErrors); ok {
		writer.Header().Add("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)

		ErrResponse := web.ResponseJSON{
			Code:   http.StatusBadRequest,
			Status: "Bad Request",
			Data:   validationError.Error(),
		}

		helpers.EncodeJSONFromResponse(writer, ErrResponse)

		return true
	}
	return false
}

func notFoundError(writer http.ResponseWriter, _ *http.Request, err interface{}) bool {

	if notFoundErr, ok := err.(NotFoundError); ok {
		writer.Header().Add("Content-Type", "application/json")
		writer.WriteHeader(http.StatusNotFound)

		ErrResponse := web.ResponseJSON{
			Code:   http.StatusNotFound,
			Status: "Not Found",
			Data:   notFoundErr.Error,
		}

		helpers.EncodeJSONFromResponse(writer, ErrResponse)

		return true
	}
	return false
}

func internalServerError(writer http.ResponseWriter, _ *http.Request, err interface{}) {
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(http.StatusInternalServerError)

	ErrResponse := web.ResponseJSON{
		Code:   http.StatusInternalServerError,
		Status: "Internal Server Error",
		Data:   err,
	}

	helpers.EncodeJSONFromResponse(writer, ErrResponse)
}
