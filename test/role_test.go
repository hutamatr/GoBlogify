package test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/hutamatr/go-blog-api/helpers"
	"github.com/hutamatr/go-blog-api/model/domain"
	"github.com/hutamatr/go-blog-api/model/web"
	repositories "github.com/hutamatr/go-blog-api/repositories/role"
	"github.com/stretchr/testify/assert"
)

func TestCreateRole(t *testing.T) {
	db := ConnectDBTest()
	DeleteDBTest(db)
	router := SetupRouterTest(db)
	defer db.Close()

	t.Run("success create role", func(t *testing.T) {
		roleBody := strings.NewReader(`{
			"name": "user"
		}`)

		request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/role", roleBody)
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
		assert.Equal(t, "user", responseBody.Data.(map[string]interface{})["name"])
	})

	t.Run("bad request create role", func(t *testing.T) {
		roleBody := strings.NewReader(`{
			"name": ""
		}`)
		request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/role", roleBody)
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

func TestFindAllRole(t *testing.T) {
	db := ConnectDBTest()
	DeleteDBTest(db)
	router := SetupRouterTest(db)
	defer db.Close()

	t.Run("success find all role", func(t *testing.T) {
		ctx := context.Background()
		tx, err := db.Begin()
		helpers.PanicError(err)

		roleRepository := repositories.NewRoleRepository()
		role1 := roleRepository.Save(ctx, tx, domain.Role{Name: "user"})
		role2 := roleRepository.Save(ctx, tx, domain.Role{Name: "admin"})

		tx.Commit()

		request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/role", nil)
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
		assert.Equal(t, role1.Name, responseBody.Data.([]interface{})[0].(map[string]interface{})["name"])
		assert.Equal(t, role2.Name, responseBody.Data.([]interface{})[1].(map[string]interface{})["name"])
	})

	t.Run("empty find all role", func(t *testing.T) {
		DeleteDBTest(db)

		request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/role", nil)
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

func TestFindByIdRole(t *testing.T) {
	db := ConnectDBTest()
	DeleteDBTest(db)
	router := SetupRouterTest(db)
	defer db.Close()

	t.Run("success find by id role", func(t *testing.T) {
		ctx := context.Background()
		tx, err := db.Begin()
		helpers.PanicError(err)

		roleRepository := repositories.NewRoleRepository()
		role := roleRepository.Save(ctx, tx, domain.Role{Name: "user"})

		tx.Commit()

		request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/role/"+strconv.Itoa(role.Id), nil)
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
		assert.Equal(t, role.Name, responseBody.Data.(map[string]interface{})["name"])
	})

	t.Run("not found find by id role", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/role/1", nil)
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
		assert.Equal(t, "role not found", responseBody.Data)
	})
}

func TestUpdateRole(t *testing.T) {
	db := ConnectDBTest()
	DeleteDBTest(db)
	router := SetupRouterTest(db)
	defer db.Close()

	t.Run("success update role", func(t *testing.T) {
		ctx := context.Background()
		tx, err := db.Begin()
		helpers.PanicError(err)

		roleRepository := repositories.NewRoleRepository()
		role := roleRepository.Save(ctx, tx, domain.Role{Name: "user"})

		tx.Commit()

		roleBody := strings.NewReader(`{
			"name": "admin"
		}`)

		request := httptest.NewRequest(http.MethodPut, "http://localhost:8080/api/role/"+strconv.Itoa(role.Id), roleBody)
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
		assert.Equal(t, "admin", responseBody.Data.(map[string]interface{})["name"])
	})

	t.Run("not found update role", func(t *testing.T) {

		roleBody := strings.NewReader(`{
			"name": "user"
		}`)

		request := httptest.NewRequest(http.MethodPut, "http://localhost:8080/api/role/1", roleBody)
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
		assert.Equal(t, "role not found", responseBody.Data)
	})

	t.Run("bad request update role", func(t *testing.T) {
		ctx := context.Background()
		tx, err := db.Begin()
		helpers.PanicError(err)

		roleRepository := repositories.NewRoleRepository()
		role := roleRepository.Save(ctx, tx, domain.Role{Name: "user"})

		tx.Commit()

		RoleBody := strings.NewReader(`{
			"name": ""
		}`)

		request := httptest.NewRequest(http.MethodPut, "http://localhost:8080/api/role/"+strconv.Itoa(role.Id), RoleBody)
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
		assert.Equal(t, "Key: 'RoleUpdateRequest.Name' Error:Field validation for 'Name' failed on the 'required' tag", responseBody.Data)
	})

}

func TestDeleteRole(t *testing.T) {
	db := ConnectDBTest()
	DeleteDBTest(db)
	router := SetupRouterTest(db)
	defer db.Close()

	t.Run("success delete role", func(t *testing.T) {
		ctx := context.Background()
		tx, err := db.Begin()
		helpers.PanicError(err)

		roleRepository := repositories.NewRoleRepository()
		role := roleRepository.Save(ctx, tx, domain.Role{Name: "Role-3"})

		tx.Commit()

		request := httptest.NewRequest(http.MethodDelete, "http://localhost:8080/api/role/"+strconv.Itoa(role.Id), nil)
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

	t.Run("not found delete Role", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodDelete, "http://localhost:8080/api/role/1", nil)
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
