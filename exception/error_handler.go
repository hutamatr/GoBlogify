package exception

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/hutamatr/GoBlogify/helpers"
)

func ErrorHandler(writer http.ResponseWriter, request *http.Request, err interface{}) {
	if validationError(writer, request, err) {
		return
	}
	if badRequestError(writer, request, err) {
		return
	}
	if notFoundError(writer, request, err) {
		return
	}
	internalServerError(writer, request, err)
}

func badRequestError(writer http.ResponseWriter, _ *http.Request, err interface{}) bool {
	exception, ok := err.(BadRequestError)
	if ok {
		writer.Header().Add("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)

		webResponseError := helpers.ErrorResponseJSON{
			Code:    http.StatusBadRequest,
			Status:  "BAD REQUEST",
			Error:   exception.Error,
			Message: "Request is not valid",
		}

		helpers.EncodeJSONFromResponse(writer, webResponseError)
		return true
	} else {
		return false
	}
}

func validationError(writer http.ResponseWriter, _ *http.Request, err interface{}) bool {
	if validationError, ok := err.(validator.ValidationErrors); ok {
		writer.Header().Add("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)

		ErrResponse := helpers.ErrorResponseJSON{
			Code:    http.StatusBadRequest,
			Status:  "BAD REQUEST",
			Error:   validationError.Error(),
			Message: "Request is not valid",
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

		ErrResponse := helpers.ErrorResponseJSON{
			Code:    http.StatusNotFound,
			Status:  "NOT FOUND",
			Error:   notFoundErr.Error,
			Message: "Resource not found",
		}

		helpers.EncodeJSONFromResponse(writer, ErrResponse)

		return true
	}
	return false
}

func internalServerError(writer http.ResponseWriter, _ *http.Request, err interface{}) {
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(http.StatusInternalServerError)

	ErrResponse := helpers.ErrorResponseJSON{
		Code:    http.StatusInternalServerError,
		Status:  "INTERNAL SERVER ERROR",
		Error:   err,
		Message: "Something went wrong",
	}

	helpers.EncodeJSONFromResponse(writer, ErrResponse)
}
