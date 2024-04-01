package test

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/hutamatr/GoBlogify/helpers"
	"github.com/hutamatr/GoBlogify/model/web"
	repositoriesRole "github.com/hutamatr/GoBlogify/repositories/role"
	repositoriesUser "github.com/hutamatr/GoBlogify/repositories/user"
	servicesAdmin "github.com/hutamatr/GoBlogify/services/admin"
	"github.com/stretchr/testify/assert"
)

func createAdminTestAdmin(db *sql.DB) (web.AdminResponse, string) {
	env := helpers.NewEnv()
	adminCode := env.Auth.AdminCode

	ctx := context.Background()
	tx, err := db.Begin()
	validator := validator.New()
	helpers.PanicError(err)
	defer tx.Commit()

	userRepository := repositoriesUser.NewUserRepository()
	roleRepository := repositoriesRole.NewRoleRepository()
	userService := servicesAdmin.NewAdminService(userRepository, roleRepository, db, validator)
	admin, accessToken, _ := userService.SignUpAdmin(ctx, web.AdminCreateRequest{Username: "admin", Email: "admin@example.com", Password: "Admin123!", Admin_Code: adminCode})

	return admin, accessToken
}

func TestCreateAdminAccount(t *testing.T) {
	env := helpers.NewEnv()
	adminCode := env.Auth.AdminCode

	db := ConnectDBTest()
	DeleteDBTest(db)
	router := SetupRouterTest(db)
	defer db.Close()

	t.Run("success create admin account", func(t *testing.T) {
		accountBody := strings.NewReader(`{
			"username": "admin",
			"email": "admin@example.com",
			"password": "Admin123!",
			"admin_code": "` + adminCode + `"
		}`)

		request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/signup-admin", accountBody)
		request.Header.Add("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, request)

		response := recorder.Result()

		assert.Equal(t, http.StatusCreated, response.StatusCode)

		body, err := io.ReadAll(response.Body)

		var responseBody web.ResponseJSON

		json.Unmarshal(body, &responseBody)

		helpers.PanicError(err)

		assert.Equal(t, http.StatusCreated, responseBody.Code)
		assert.Equal(t, "CREATED", responseBody.Status)
		assert.Equal(t, "admin", responseBody.Data.(map[string]interface{})["admin"].(map[string]interface{})["username"])
		assert.Equal(t, "admin@example.com", responseBody.Data.(map[string]interface{})["admin"].(map[string]interface{})["email"])
	})

	t.Run("failed create admin account", func(t *testing.T) {
		accountBody := strings.NewReader(`{
			"username": "",
			"email": "admin@example.com",
			"password": "Admin123!",
			"admin_code": "` + adminCode + `"
		}`)

		request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/signup-admin", accountBody)
		request.Header.Add("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, request)

		response := recorder.Result()

		assert.Equal(t, http.StatusBadRequest, response.StatusCode)

		body, err := io.ReadAll(response.Body)

		var responseBody web.ResponseJSON

		json.Unmarshal(body, &responseBody)

		helpers.PanicError(err)

		assert.Equal(t, http.StatusBadRequest, responseBody.Code)
		assert.Equal(t, "Bad Request", responseBody.Status)
	})
}

func TestLoginAdmin(t *testing.T) {
	db := ConnectDBTest()
	DeleteDBTest(db)
	router := SetupRouterTest(db)
	defer db.Close()

	admin, _ := createAdminTestAdmin(db)

	t.Run("success login admin", func(t *testing.T) {
		accountBody := strings.NewReader(`{
			"email": "admin@example.com",
			"password": "Admin123!"
		}`)

		request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/signin-admin", accountBody)
		request.Header.Add("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, request)

		response := recorder.Result()

		assert.Equal(t, http.StatusOK, response.StatusCode)

		body, err := io.ReadAll(response.Body)

		var responseBody web.ResponseJSON

		json.Unmarshal(body, &responseBody)

		helpers.PanicError(err)

		assert.Equal(t, http.StatusOK, responseBody.Code)
		assert.Equal(t, "OK", responseBody.Status)
		assert.Equal(t, admin.Email, responseBody.Data.(map[string]interface{})["admin"].(map[string]interface{})["email"])
		assert.Equal(t, admin.Username, responseBody.Data.(map[string]interface{})["admin"].(map[string]interface{})["username"])
	})

	t.Run("failed login admin", func(t *testing.T) {
		accountBody := strings.NewReader(`{
			"email": "admin@example.com",
			"password": ""
		}`)

		request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/signin-admin", accountBody)
		request.Header.Add("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, request)

		response := recorder.Result()

		assert.Equal(t, http.StatusBadRequest, response.StatusCode)

		body, err := io.ReadAll(response.Body)

		var responseBody web.ResponseJSON

		json.Unmarshal(body, &responseBody)

		helpers.PanicError(err)

		assert.Equal(t, http.StatusBadRequest, responseBody.Code)
		assert.Equal(t, "Bad Request", responseBody.Status)
	})
}
