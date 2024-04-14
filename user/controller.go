package user

import (
	"net/http"
	"strconv"
	"time"

	"github.com/hutamatr/GoBlogify/exception"
	"github.com/hutamatr/GoBlogify/helpers"
	"github.com/julienschmidt/httprouter"
)

type UserController interface {
	CreateUserHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	SignInUserHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	SignOutUserHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindAllUserHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindByIdUserHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	UpdateUserHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	DeleteUserHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	GetRefreshTokenHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}

type UserControllerImpl struct {
	service UserService
}

func NewUserController(userService UserService) UserController {
	return &UserControllerImpl{
		service: userService,
	}
}

func (controller *UserControllerImpl) CreateUserHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	var env = helpers.NewEnv()
	var AppEnv = env.App.AppEnv

	var newUserRequest UserCreateRequest

	helpers.DecodeJSONFromRequest(request, &newUserRequest)

	newUser, accessToken, refreshToken := controller.service.SignUp(request.Context(), newUserRequest)

	cookie := http.Cookie{}
	cookie.Name = "rt"
	cookie.Value = refreshToken
	cookie.MaxAge = 7 * 24 * 60 * 60
	cookie.Secure = AppEnv == "production"
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteStrictMode
	cookie.Expires = time.Now().Add(7 * 24 * time.Hour)
	http.SetCookie(writer, &cookie)

	userResponse := helpers.ResponseJSON{
		Code:   http.StatusCreated,
		Status: "CREATED",
		Data: map[string]interface{}{
			"access_token": accessToken,
			"user":         newUser,
		},
	}

	writer.WriteHeader(http.StatusCreated)
	helpers.EncodeJSONFromResponse(writer, userResponse)
}

func (controller *UserControllerImpl) SignInUserHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	var env = helpers.NewEnv()
	var AppEnv = env.App.AppEnv

	var signInRequest UserLoginRequest

	helpers.DecodeJSONFromRequest(request, &signInRequest)

	signInUser, accessToken, refreshToken := controller.service.SignIn(request.Context(), signInRequest)

	cookie := http.Cookie{}
	cookie.Name = "rt"
	cookie.Value = refreshToken
	cookie.MaxAge = 7 * 24 * 60 * 60
	cookie.Secure = AppEnv == "production"
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteStrictMode
	cookie.Expires = time.Now().Add(7 * 24 * time.Hour)
	http.SetCookie(writer, &cookie)

	userResponse := helpers.ResponseJSON{
		Code:   http.StatusOK,
		Status: "OK",
		Data: map[string]interface{}{
			"access_token": accessToken,
			"user":         signInUser,
		},
	}

	writer.WriteHeader(http.StatusOK)
	helpers.EncodeJSONFromResponse(writer, userResponse)
}

func (controller *UserControllerImpl) SignOutUserHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	var env = helpers.NewEnv()
	var AppEnv = env.App.AppEnv

	cookie := http.Cookie{}
	cookie.Name = "rt"
	cookie.Value = ""
	cookie.MaxAge = -1
	cookie.Secure = AppEnv == "production"
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteStrictMode
	cookie.Expires = time.Now().Add(7 * 24 * time.Hour)
	http.SetCookie(writer, &cookie)

	userResponse := helpers.ResponseJSON{
		Code:   http.StatusOK,
		Status: "OK",
	}

	writer.WriteHeader(http.StatusOK)
	helpers.EncodeJSONFromResponse(writer, userResponse)
}

func (controller *UserControllerImpl) FindAllUserHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	isAdmin := helpers.IsAdmin(request)

	users := controller.service.FindAll(request.Context(), isAdmin)

	userResponse := helpers.ResponseJSON{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   users,
	}

	writer.WriteHeader(http.StatusOK)
	helpers.EncodeJSONFromResponse(writer, userResponse)
}

func (controller *UserControllerImpl) FindByIdUserHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id := params.ByName("userId")
	userId, err := strconv.Atoi(id)
	helpers.PanicError(err, "Invalid User Id")

	user := controller.service.FindById(request.Context(), userId)

	userResponse := helpers.ResponseJSON{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   user,
	}

	writer.WriteHeader(http.StatusOK)
	helpers.EncodeJSONFromResponse(writer, userResponse)
}

func (controller *UserControllerImpl) UpdateUserHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id := params.ByName("userId")
	userId, err := strconv.Atoi(id)
	helpers.PanicError(err, "Invalid User Id")

	var updateUserRequest UserUpdateRequest

	updateUserRequest.Id = userId

	helpers.DecodeJSONFromRequest(request, &updateUserRequest)

	updatedUser := controller.service.Update(request.Context(), updateUserRequest)

	userResponse := helpers.ResponseJSON{
		Code:   http.StatusOK,
		Status: "UPDATED",
		Data:   updatedUser,
	}

	writer.WriteHeader(http.StatusOK)
	helpers.EncodeJSONFromResponse(writer, userResponse)
}

func (controller *UserControllerImpl) DeleteUserHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id := params.ByName("userId")
	userId, err := strconv.Atoi(id)
	helpers.PanicError(err, "Invalid User Id")

	controller.service.Delete(request.Context(), userId)

	userResponse := helpers.ResponseJSON{
		Code:   http.StatusOK,
		Status: "DELETED",
	}

	writer.WriteHeader(http.StatusOK)
	helpers.EncodeJSONFromResponse(writer, userResponse)
}

func (controller *UserControllerImpl) GetRefreshTokenHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	var env = helpers.NewEnv()
	var AppEnv = env.App.AppEnv
	var accessTokenSecret = env.SecretToken.AccessSecret
	var refreshTokenSecret = env.SecretToken.RefreshSecret

	refreshToken, err := request.Cookie("rt")
	if err != nil {
		panic(exception.NewUnauthorizedError(err.Error()))
	}

	claims, err := helpers.VerifyToken(refreshToken.Value, []byte(refreshTokenSecret))

	if err != nil {
		panic(exception.NewUnauthorizedError(err.Error()))
	}

	idFloat := claims["sub"].(float64)
	userId := int(idFloat)

	accessTokenExpired := helpers.AccessTokenDuration(AppEnv)

	newAccessToken, err := helpers.GenerateToken(userId, accessTokenExpired, accessTokenSecret)

	helpers.PanicError(err, "failed to generate access token")

	userResponse := helpers.ResponseJSON{
		Code:   http.StatusOK,
		Status: "OK",
		Data: map[string]interface{}{
			"access_token": newAccessToken,
		},
	}

	writer.WriteHeader(http.StatusOK)
	helpers.EncodeJSONFromResponse(writer, userResponse)
}
