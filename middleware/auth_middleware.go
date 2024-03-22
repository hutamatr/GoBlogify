package middleware

import (
	"net/http"
	"strings"

	"github.com/hutamatr/GoBlogify/helpers"
	"github.com/hutamatr/GoBlogify/model/web"
)

type AuthMiddleware struct {
	Handler http.Handler
}

var publicRoutes = []string{
	"/api/signup",
	"/api/signin",
	"/api/signout",
	"/api/refresh",
}

func NewAuthMiddleware(handler http.Handler) *AuthMiddleware {
	return &AuthMiddleware{Handler: handler}
}

var env = helpers.NewEnv()
var tokenSecret = env.SecretToken.AccessSecret

func (middleware *AuthMiddleware) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	path := request.URL.Path

	for _, publicRoute := range publicRoutes {
		if publicRoute == path {
			middleware.Handler.ServeHTTP(writer, request)
			return
		}
	}

	authorizationHeader := request.Header.Get("Authorization")

	if tokenString := strings.TrimSpace(strings.Replace(authorizationHeader, "Bearer ", "", 1)); tokenString == "" {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusUnauthorized)

		ErrResponse := web.ResponseJSON{
			Code:   http.StatusUnauthorized,
			Status: "Unauthorized",
			Data:   "token is required",
		}

		helpers.EncodeJSONFromResponse(writer, ErrResponse)
		return
	} else if _, err := helpers.VerifyToken(tokenString, []byte(tokenSecret)); err != nil {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusUnauthorized)

		ErrResponse := web.ResponseJSON{
			Code:   http.StatusUnauthorized,
			Status: "Unauthorized",
			Data:   err.Error(),
		}

		helpers.EncodeJSONFromResponse(writer, ErrResponse)
		return
	} else {
		middleware.Handler.ServeHTTP(writer, request)
	}
}
