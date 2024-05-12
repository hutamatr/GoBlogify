package admin

import (
	"net/http"
	"time"

	"github.com/hutamatr/GoBlogify/helpers"
	"github.com/julienschmidt/httprouter"
)

type AdminController interface {
	CreateAdminHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	SignInAdminHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}

type AdminControllerImpl struct {
	service AdminService
}

func NewAdminController(service AdminService) AdminController {
	return &AdminControllerImpl{
		service: service,
	}
}

func (controller *AdminControllerImpl) CreateAdminHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	var env = helpers.NewEnv()
	var AppEnv = env.App.AppEnv
	var newAdminRequest AdminCreateRequest

	helpers.DecodeJSONFromRequest(request, &newAdminRequest)

	newAdmin, accessToken, refreshToken := controller.service.SignUpAdmin(request.Context(), newAdminRequest)

	cookie := http.Cookie{}
	cookie.Name = "rt"
	cookie.Value = refreshToken
	cookie.MaxAge = 7 * 24 * 60 * 60
	cookie.Secure = AppEnv == "production"
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteStrictMode
	cookie.Expires = time.Now().Add(7 * 24 * time.Hour)
	http.SetCookie(writer, &cookie)

	adminResponse := helpers.ResponseJSON{
		Code:   http.StatusCreated,
		Status: "CREATED",
		Data: map[string]interface{}{
			"access_token": accessToken,
			"admin":        newAdmin,
		},
	}
	writer.WriteHeader(http.StatusCreated)

	helpers.EncodeJSONFromResponse(writer, adminResponse)
}

func (controller *AdminControllerImpl) SignInAdminHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	var env = helpers.NewEnv()
	var AppEnv = env.App.AppEnv
	var signInRequest AdminLoginRequest

	helpers.DecodeJSONFromRequest(request, &signInRequest)

	signInAdmin, accessToken, refreshToken := controller.service.SignInAdmin(request.Context(), signInRequest)

	cookie := http.Cookie{}
	cookie.Name = "rt"
	cookie.Value = refreshToken
	cookie.MaxAge = 7 * 24 * 60 * 60
	cookie.Secure = AppEnv == "production"
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteStrictMode
	cookie.Expires = time.Now().Add(7 * 24 * time.Hour)
	http.SetCookie(writer, &cookie)

	adminResponse := helpers.ResponseJSON{
		Code:   http.StatusOK,
		Status: "OK",
		Data: map[string]interface{}{
			"access_token": accessToken,
			"admin":        signInAdmin,
		},
	}

	writer.WriteHeader(http.StatusOK)
	helpers.EncodeJSONFromResponse(writer, adminResponse)
}
