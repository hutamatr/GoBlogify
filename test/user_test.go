package test

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/hutamatr/GoBlogify/helpers"
	"github.com/hutamatr/GoBlogify/model/web"
	repositoriesRole "github.com/hutamatr/GoBlogify/repositories/role"
	repositoriesUser "github.com/hutamatr/GoBlogify/repositories/user"
	servicesUser "github.com/hutamatr/GoBlogify/services/user"

	"github.com/stretchr/testify/assert"
)

func createUserTestUser(db *sql.DB) web.UserResponse {
	ctx := context.Background()
	tx, err := db.Begin()
	validator := validator.New()
	helpers.PanicError(err)
	defer tx.Commit()

	userRepository := repositoriesUser.NewUserRepository()
	roleRepository := repositoriesRole.NewRoleRepository()
	userService := servicesUser.NewUserService(userRepository, roleRepository, db, validator)
	user, _, _ := userService.SignUp(ctx, web.UserCreateRequest{Username: "userTest", Email: "testing@example.com", Password: "Password123!"})

	return user
}

func TestCreateAccount(t *testing.T) {
	db := ConnectDBTest()
	DeleteDBTest(db)
	router := SetupRouterTest(db)
	defer db.Close()

	t.Run("success create account", func(t *testing.T) {
		accountBody := strings.NewReader(`{
			"username": "userTest",
			"email": "testing@example.com",
			"password": "Password123!"
		}`)

		request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/signup", accountBody)
		request.Header.Add("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, request)

		response := recorder.Result()

		assert.Equal(t, http.StatusOK, response.StatusCode)

		body, err := io.ReadAll(response.Body)

		var responseBody web.ResponseJSON

		json.Unmarshal(body, &responseBody)

		helpers.PanicError(err)

		assert.Equal(t, http.StatusCreated, responseBody.Code)
		assert.Equal(t, "CREATED", responseBody.Status)
		assert.Equal(t, "userTest", responseBody.Data.(map[string]interface{})["user"].(map[string]interface{})["username"])
		assert.Equal(t, "testing@example.com", responseBody.Data.(map[string]interface{})["user"].(map[string]interface{})["email"])
	})

	t.Run("failed create account", func(t *testing.T) {
		accountBody := strings.NewReader(`{
			"username": "",
			"email": "testing@example.com",
			"password": "Password123!"
		}`)

		request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/signup", accountBody)
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

func TestLogin(t *testing.T) {
	db := ConnectDBTest()
	DeleteDBTest(db)
	router := SetupRouterTest(db)
	defer db.Close()

	user := createUserTestUser(db)

	t.Run("success login", func(t *testing.T) {
		accountBody := strings.NewReader(`{
			"email": "testing@example.com",
			"password": "Password123!"
		}`)

		request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/signin", accountBody)
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
		assert.Equal(t, user.Email, responseBody.Data.(map[string]interface{})["user"].(map[string]interface{})["email"])
		assert.Equal(t, user.Username, responseBody.Data.(map[string]interface{})["user"].(map[string]interface{})["username"])
	})

	t.Run("failed login", func(t *testing.T) {
		accountBody := strings.NewReader(`{
			"email": "testing@example.com",
			"password": ""
		}`)

		request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/signin", accountBody)
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

func TestFindAllUser(t *testing.T) {
	db := ConnectDBTest()
	DeleteDBTest(db)
	router := SetupRouterTest(db)
	defer db.Close()

	createUserTestUser(db)

	t.Run("success find all user", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/user", nil)
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
	})

	t.Run("failed find all user", func(t *testing.T) {
		DeleteDBTest(db)
		request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/user", nil)
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
		assert.Nil(t, responseBody.Data)
	})
}

func TestFindByIdUser(t *testing.T) {
	db := ConnectDBTest()
	DeleteDBTest(db)
	router := SetupRouterTest(db)
	defer db.Close()

	user := createUserTestUser(db)

	t.Run("success find by id user", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/user/"+strconv.Itoa(user.Id), nil)
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
		assert.Equal(t, user.Email, responseBody.Data.(map[string]interface{})["email"])
		assert.Equal(t, user.Username, responseBody.Data.(map[string]interface{})["username"])
	})

	t.Run("failed find by id user", func(t *testing.T) {
		DeleteDBTest(db)
		request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/user/0", nil)
		request.Header.Add("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, request)

		response := recorder.Result()

		assert.Equal(t, http.StatusNotFound, response.StatusCode)

		body, err := io.ReadAll(response.Body)

		var responseBody web.ResponseJSON

		json.Unmarshal(body, &responseBody)

		helpers.PanicError(err)

		assert.Equal(t, http.StatusNotFound, responseBody.Code)
		assert.Equal(t, "Not Found", responseBody.Status)
	})
}

func TestUpdateUser(t *testing.T) {
	db := ConnectDBTest()
	DeleteDBTest(db)
	router := SetupRouterTest(db)
	defer db.Close()

	user := createUserTestUser(db)

	t.Run("success update user", func(t *testing.T) {
		accountBody := strings.NewReader(`{
			"role_id" : ` + strconv.Itoa(user.Role_Id) + `,
			"username": "userTest2",
			"first_name" : "test",
			"last_name" : "testing" 
		}`)

		request := httptest.NewRequest(http.MethodPut, "http://localhost:8080/api/user/"+strconv.Itoa(user.Id), accountBody)
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
		assert.Equal(t, "UPDATED", responseBody.Status)
		assert.Equal(t, user.Email, responseBody.Data.(map[string]interface{})["email"])
		assert.Equal(t, "userTest2", responseBody.Data.(map[string]interface{})["username"])
	})

	t.Run("failed update user", func(t *testing.T) {
		accountBody := strings.NewReader(`{
			"role_id" : ` + strconv.Itoa(user.Role_Id) + `,
			"username": "",
			"first_name" : "",
			"last_name" : "" 
		}`)

		request := httptest.NewRequest(http.MethodPut, "http://localhost:8080/api/user/"+strconv.Itoa(user.Id), accountBody)
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

func TestDeleteUser(t *testing.T) {
	db := ConnectDBTest()
	DeleteDBTest(db)
	router := SetupRouterTest(db)
	defer db.Close()

	user := createUserTestUser(db)

	t.Run("success delete user", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodDelete, "http://localhost:8080/api/user/"+strconv.Itoa(user.Id), nil)
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
		assert.Equal(t, "DELETED", responseBody.Status)
	})

	t.Run("failed delete user", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodDelete, "http://localhost:8080/api/user/0", nil)
		request.Header.Add("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, request)

		response := recorder.Result()

		assert.Equal(t, http.StatusNotFound, response.StatusCode)

		body, err := io.ReadAll(response.Body)

		var responseBody web.ResponseJSON

		json.Unmarshal(body, &responseBody)

		helpers.PanicError(err)

		assert.Equal(t, http.StatusNotFound, responseBody.Code)
		assert.Equal(t, "Not Found", responseBody.Status)
	})
}
