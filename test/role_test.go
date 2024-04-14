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

	"github.com/hutamatr/GoBlogify/helpers"
	"github.com/hutamatr/GoBlogify/role"
	"github.com/stretchr/testify/assert"
)

func TestCreateRole(t *testing.T) {
	db := ConnectDBTest()
	DeleteDBTest(db)
	router := SetupRouterTest(db)
	defer db.Close()

	_, accessToken := createAdminTestAdmin(db)

	t.Run("success create role", func(t *testing.T) {
		roleBody := strings.NewReader(`{
			"name": "user"
		}`)

		request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/role", roleBody)
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", "Bearer "+accessToken)

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
		assert.Equal(t, "user", responseBody.Data.(map[string]interface{})["name"])
	})

	t.Run("bad request create role", func(t *testing.T) {
		roleBody := strings.NewReader(`{
			"name": ""
		}`)
		request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/role", roleBody)
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", "Bearer "+accessToken)

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, request)

		response := recorder.Result()

		assert.Equal(t, http.StatusBadRequest, response.StatusCode)

		body, err := io.ReadAll(response.Body)

		var responseBody helpers.ResponseJSON

		json.Unmarshal(body, &responseBody)

		helpers.PanicError(err, "failed to read response body")

		assert.Equal(t, http.StatusBadRequest, responseBody.Code)
		assert.Equal(t, "BAD REQUEST", responseBody.Status)
	})
}

func TestFindAllRole(t *testing.T) {
	db := ConnectDBTest()
	DeleteDBTest(db)
	router := SetupRouterTest(db)
	defer db.Close()

	_, accessToken := createAdminTestAdmin(db)

	t.Run("success find all role", func(t *testing.T) {
		ctx := context.Background()
		tx, err := db.Begin()
		helpers.PanicError(err, "failed to begin transaction")

		roleRepository := role.NewRoleRepository()
		role1 := roleRepository.Save(ctx, tx, role.Role{Name: "user"})

		tx.Commit()

		request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/role", nil)
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", "Bearer "+accessToken)

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
		assert.Equal(t, "admin", responseBody.Data.([]interface{})[0].(map[string]interface{})["name"])
		assert.Equal(t, role1.Name, responseBody.Data.([]interface{})[1].(map[string]interface{})["name"])
	})

	t.Run("empty find all role", func(t *testing.T) {
		DeleteDBTest(db)

		request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/role", nil)
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", "Bearer "+accessToken)

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
		assert.Nil(t, responseBody.Data)
	})
}

func TestFindByIdRole(t *testing.T) {
	db := ConnectDBTest()
	DeleteDBTest(db)
	router := SetupRouterTest(db)
	defer db.Close()

	_, accessToken := createAdminTestAdmin(db)

	t.Run("success find by id role", func(t *testing.T) {
		ctx := context.Background()
		tx, err := db.Begin()
		helpers.PanicError(err, "failed to begin transaction")

		roleRepository := role.NewRoleRepository()
		role := roleRepository.Save(ctx, tx, role.Role{Name: "user"})

		tx.Commit()

		request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/role/"+strconv.Itoa(role.Id), nil)
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", "Bearer "+accessToken)

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
		assert.Equal(t, role.Name, responseBody.Data.(map[string]interface{})["name"])
	})

	t.Run("not found find by id role", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/role/10", nil)
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", "Bearer "+accessToken)

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, request)

		response := recorder.Result()

		assert.Equal(t, http.StatusNotFound, response.StatusCode)

		body, err := io.ReadAll(response.Body)

		var responseBody helpers.ResponseJSON

		json.Unmarshal(body, &responseBody)

		helpers.PanicError(err, "failed to read response body")

		assert.Equal(t, http.StatusNotFound, responseBody.Code)
		assert.Equal(t, "NOT FOUND", responseBody.Status)
		assert.Equal(t, "role not found", responseBody.Data)
	})
}

func TestUpdateRole(t *testing.T) {
	db := ConnectDBTest()
	DeleteDBTest(db)
	router := SetupRouterTest(db)
	defer db.Close()

	_, accessToken := createAdminTestAdmin(db)

	t.Run("success update role", func(t *testing.T) {
		ctx := context.Background()
		tx, err := db.Begin()
		helpers.PanicError(err, "failed to begin transaction")

		roleRepository := role.NewRoleRepository()
		role := roleRepository.Save(ctx, tx, role.Role{Name: "user"})

		tx.Commit()

		roleBody := strings.NewReader(`{
			"name": "user-2"
		}`)

		request := httptest.NewRequest(http.MethodPut, "http://localhost:8080/api/role/"+strconv.Itoa(role.Id), roleBody)
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", "Bearer "+accessToken)

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, request)

		response := recorder.Result()

		assert.Equal(t, http.StatusOK, response.StatusCode)

		body, err := io.ReadAll(response.Body)

		var responseBody helpers.ResponseJSON

		json.Unmarshal(body, &responseBody)

		helpers.PanicError(err, "failed to read response body")

		assert.Equal(t, http.StatusOK, responseBody.Code)
		assert.Equal(t, "UPDATED", responseBody.Status)
		assert.Equal(t, "user-2", responseBody.Data.(map[string]interface{})["name"])
	})

	t.Run("not found update role", func(t *testing.T) {

		roleBody := strings.NewReader(`{
			"name": "user"
		}`)

		request := httptest.NewRequest(http.MethodPut, "http://localhost:8080/api/role/10", roleBody)
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", "Bearer "+accessToken)

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, request)

		response := recorder.Result()

		assert.Equal(t, http.StatusNotFound, response.StatusCode)

		body, err := io.ReadAll(response.Body)

		var responseBody helpers.ResponseJSON

		json.Unmarshal(body, &responseBody)

		helpers.PanicError(err, "failed to read response body")

		assert.Equal(t, http.StatusNotFound, responseBody.Code)
		assert.Equal(t, "NOT FOUND", responseBody.Status)
		assert.Equal(t, "role not found", responseBody.Data)
	})

	t.Run("bad request update role", func(t *testing.T) {
		ctx := context.Background()
		tx, err := db.Begin()
		helpers.PanicError(err, "failed to begin transaction")

		roleRepository := role.NewRoleRepository()
		role := roleRepository.Save(ctx, tx, role.Role{Name: "user"})

		tx.Commit()

		RoleBody := strings.NewReader(`{
			"name": ""
		}`)

		request := httptest.NewRequest(http.MethodPut, "http://localhost:8080/api/role/"+strconv.Itoa(role.Id), RoleBody)
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", "Bearer "+accessToken)

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, request)

		response := recorder.Result()

		assert.Equal(t, http.StatusBadRequest, response.StatusCode)

		body, err := io.ReadAll(response.Body)

		var responseBody helpers.ResponseJSON

		json.Unmarshal(body, &responseBody)

		helpers.PanicError(err, "failed to read response body")

		assert.Equal(t, http.StatusBadRequest, responseBody.Code)
		assert.Equal(t, "BAD REQUEST", responseBody.Status)
		assert.Equal(t, "Key: 'RoleUpdateRequest.Name' Error:Field validation for 'Name' failed on the 'required' tag", responseBody.Data)
	})

}

func TestDeleteRole(t *testing.T) {
	db := ConnectDBTest()
	DeleteDBTest(db)
	router := SetupRouterTest(db)
	defer db.Close()

	_, accessToken := createAdminTestAdmin(db)

	t.Run("success delete role", func(t *testing.T) {
		ctx := context.Background()
		tx, err := db.Begin()
		helpers.PanicError(err, "failed to begin transaction")

		roleRepository := role.NewRoleRepository()
		role := roleRepository.Save(ctx, tx, role.Role{Name: "Role-3"})

		tx.Commit()

		request := httptest.NewRequest(http.MethodDelete, "http://localhost:8080/api/role/"+strconv.Itoa(role.Id), nil)
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", "Bearer "+accessToken)

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, request)

		response := recorder.Result()

		assert.Equal(t, http.StatusOK, response.StatusCode)

		body, err := io.ReadAll(response.Body)

		var responseBody helpers.ResponseJSON

		json.Unmarshal(body, &responseBody)

		helpers.PanicError(err, "failed to read response body")

		assert.Equal(t, http.StatusOK, responseBody.Code)
		assert.Equal(t, "DELETED", responseBody.Status)
	})

	t.Run("not found delete Role", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodDelete, "http://localhost:8080/api/role/10", nil)
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", "Bearer "+accessToken)

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, request)

		response := recorder.Result()

		assert.Equal(t, http.StatusNotFound, response.StatusCode)

		body, err := io.ReadAll(response.Body)

		var responseBody helpers.ResponseJSON

		json.Unmarshal(body, &responseBody)

		helpers.PanicError(err, "failed to read response body")

		assert.Equal(t, http.StatusNotFound, responseBody.Code)
		assert.Equal(t, "NOT FOUND", responseBody.Status)
	})
}
