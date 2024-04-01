package controllers

import (
	"net/http"
	"time"

	"github.com/hutamatr/GoBlogify/helpers"
	"github.com/hutamatr/GoBlogify/model/web"
	servicesAdmin "github.com/hutamatr/GoBlogify/services/admin"
	"github.com/julienschmidt/httprouter"
)

type AdminControllerImpl struct {
	service servicesAdmin.AdminService
}

func NewAdminControllerImpl(service servicesAdmin.AdminService) AdminController {
	return &AdminControllerImpl{
		service: service,
	}
}

func (controller *AdminControllerImpl) CreateAdmin(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	var env = helpers.NewEnv()
	var AppEnv = env.App.AppEnv
	var newAdminRequest web.AdminCreateRequest

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

	adminResponse := web.ResponseJSON{
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

func (controller *AdminControllerImpl) SignInAdmin(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	var env = helpers.NewEnv()
	var AppEnv = env.App.AppEnv
	var signInRequest web.AdminLoginRequest

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

	adminResponse := web.ResponseJSON{
		Code:   http.StatusOK,
		Status: "OK",
		Data: map[string]interface{}{
			"access_token": accessToken,
			"admin":        signInAdmin,
		},
	}

	helpers.EncodeJSONFromResponse(writer, adminResponse)
}
