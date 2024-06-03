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

	"github.com/hutamatr/GoBlogify/admin"
	"github.com/hutamatr/GoBlogify/helpers"
	"github.com/hutamatr/GoBlogify/role"
	"github.com/hutamatr/GoBlogify/user"
	"github.com/stretchr/testify/assert"
)

func createAdminTestAdmin(db *sql.DB) (admin.AdminResponse, string) {
	env := helpers.NewEnv()
	adminCode := env.Auth.AdminCode
	ctx := context.Background()
	tx, err := db.Begin()
	helpers.PanicError(err, "failed to begin transaction")
	defer tx.Commit()

	userRepository := user.NewUserRepository()
	roleRepository := role.NewRoleRepository()
	userService := admin.NewAdminService(userRepository, roleRepository, db, helpers.Validate)
	admin, accessToken, _ := userService.SignUpAdmin(ctx, admin.AdminCreateRequest{Username: "admin", Email: "admin@example.com", Password: "Admin123!", Admin_Code: adminCode, Confirm_Password: "Admin123!"})

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
			"confirm_password": "Admin123!",
			"admin_code": "` + adminCode + `"
		}`)

		request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/signup-admin", accountBody)
		request.Header.Add("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, request)

		response := recorder.Result()

		assert.Equal(t, http.StatusCreated, response.StatusCode)

		body, err := io.ReadAll(response.Body)

		var responseBody helpers.ResponseJSON

		json.Unmarshal(body, &responseBody)

		helpers.PanicError(err, "failed to read response body")

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

		request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/signup-admin", accountBody)
		request.Header.Add("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, request)

		response := recorder.Result()

		assert.Equal(t, http.StatusBadRequest, response.StatusCode)

		body, err := io.ReadAll(response.Body)

		var responseBody helpers.ErrorResponseJSON

		json.Unmarshal(body, &responseBody)

		helpers.PanicError(err, "failed to read response body")

		assert.Equal(t, http.StatusBadRequest, responseBody.Code)
		assert.Equal(t, "BAD REQUEST", responseBody.Status)
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
		

		request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/signin-admin", accountBody)
		request.Header.Add("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, request)

		response := recorder.Result()

		assert.Equal(t, http.StatusOK, response.StatusCode)

		body, err := io.ReadAll(response.Body)

		var responseBody helpers.ResponseJSON

		json.Unmarshal(body, &responseBody)

		helpers.PanicError(err, "failed to read response body")

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

		request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/signin-admin", accountBody)
		request.Header.Add("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, request)

		response := recorder.Result()

		assert.Equal(t, http.StatusBadRequest, response.StatusCode)

		body, err := io.ReadAll(response.Body)

		var responseBody helpers.ErrorResponseJSON

		json.Unmarshal(body, &responseBody)

		helpers.PanicError(err, "failed to read response body")

		assert.Equal(t, http.StatusBadRequest, responseBody.Code)
		assert.Equal(t, "BAD REQUEST", responseBody.Status)
	})
}
