package middleware

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/hutamatr/GoBlogify/app"
	"github.com/hutamatr/GoBlogify/helpers"
)

type AuthMiddleware struct {
	Handler http.Handler
}

var publicRoutes = []string{
	"/api/signup",
	"/api/signin",
	"/api/signup-admin",
	"/api/signin-admin",
	"/api/signout",
	"/api/refresh",
}

func NewAuthMiddleware(handler http.Handler) *AuthMiddleware {
	return &AuthMiddleware{Handler: handler}
}

func (middleware *AuthMiddleware) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	env := helpers.NewEnv()
	tokenSecret := env.SecretToken.AccessSecret
	path := request.URL.Path

	for _, publicRoute := range publicRoutes {
		if publicRoute == path {
			middleware.Handler.ServeHTTP(writer, request)
			return
		}
	}

	authorizationHeader := request.Header.Get("Authorization")

	tokenString := strings.TrimSpace(strings.Replace(authorizationHeader, "Bearer ", "", 1))

	if tokenString == "" {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusUnauthorized)

		ErrResponse := helpers.ResponseJSON{
			Code:   http.StatusUnauthorized,
			Status: "Unauthorized",
			Data:   "token is required",
		}

		helpers.EncodeJSONFromResponse(writer, ErrResponse)
		return
	}

	claims, err := helpers.VerifyToken(tokenString, []byte(tokenSecret))

	if err != nil {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusUnauthorized)

		ErrResponse := helpers.ResponseJSON{
			Code:   http.StatusUnauthorized,
			Status: "Unauthorized",
			Data:   err.Error(),
		}

		helpers.EncodeJSONFromResponse(writer, ErrResponse)
		return
	}

	idFloat := claims["sub"].(float64)
	id := int(idFloat)

	db := app.ConnectDB()
	defer db.Close()

	queryUserRole := "SELECT role_id FROM user WHERE id = ?"
	rows, err := db.Query(queryUserRole, id)
	helpers.PanicError(err, "failed to query user role")

	defer rows.Close()
	var userRoleId int

	if rows.Next() {
		err = rows.Scan(&userRoleId)
		helpers.PanicError(err, "failed to scan user role")
	}

	queryRole := "SELECT id FROM role WHERE name = ?"
	rows2, err := db.Query(queryRole, "admin")
	helpers.PanicError(err, "failed to query role")

	defer rows2.Close()
	var roleId int
	var nullRoleId sql.NullInt32

	if rows2.Next() {
		err = rows2.Scan(&nullRoleId)
		helpers.PanicError(err, "failed to scan role")
	}

	if nullRoleId.Valid {
		roleId = int(nullRoleId.Int32)
	} else {
		roleId = 0
	}

	isAdmin := "false"
	if userRoleId == roleId {
		isAdmin = "true"
	}
	request.Header.Set("isAdmin", isAdmin)

	middleware.Handler.ServeHTTP(writer, request)
}
