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
	"github.com/hutamatr/GoBlogify/model/domain"
	"github.com/hutamatr/GoBlogify/model/web"
	repositoriesCategory "github.com/hutamatr/GoBlogify/repositories/category"
	"github.com/stretchr/testify/assert"
)

func TestCreateCategory(t *testing.T) {
	db := ConnectDBTest()
	DeleteDBTest(db)
	router := SetupRouterTest(db)
	defer db.Close()

	t.Run("success create category", func(t *testing.T) {
		categoryBody := strings.NewReader(`{
			"name": "category-1"
		}`)

		request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/category", categoryBody)
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
		assert.Equal(t, "category-1", responseBody.Data.(map[string]interface{})["name"])
	})

	t.Run("bad request create category", func(t *testing.T) {
		categoryBody := strings.NewReader(`{
			"name": ""
		}`)
		request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/category", categoryBody)
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

func TestFindAllCategory(t *testing.T) {
	db := ConnectDBTest()
	DeleteDBTest(db)
	router := SetupRouterTest(db)
	defer db.Close()

	t.Run("success find all category", func(t *testing.T) {
		ctx := context.Background()
		tx, err := db.Begin()
		helpers.PanicError(err)

		categoryRepository := repositoriesCategory.NewCategoryRepository()
		category1 := categoryRepository.Save(ctx, tx, domain.Category{Name: "category-3"})
		category2 := categoryRepository.Save(ctx, tx, domain.Category{Name: "category-4"})

		tx.Commit()

		request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/category", nil)
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
		assert.Equal(t, category1.Name, responseBody.Data.([]interface{})[0].(map[string]interface{})["name"])
		assert.Equal(t, category2.Name, responseBody.Data.([]interface{})[1].(map[string]interface{})["name"])
	})

	t.Run("empty find all category", func(t *testing.T) {
		DeleteDBTest(db)

		request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/category", nil)
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

func TestFindByIdCategory(t *testing.T) {
	db := ConnectDBTest()
	DeleteDBTest(db)
	router := SetupRouterTest(db)
	defer db.Close()

	t.Run("success find by id category", func(t *testing.T) {
		ctx := context.Background()
		tx, err := db.Begin()
		helpers.PanicError(err)

		categoryRepository := repositoriesCategory.NewCategoryRepository()
		category := categoryRepository.Save(ctx, tx, domain.Category{Name: "category-5"})

		tx.Commit()

		request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/category/"+strconv.Itoa(category.Id), nil)
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
		assert.Equal(t, category.Name, responseBody.Data.(map[string]interface{})["name"])
	})

	t.Run("not found find by id category", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/category/1", nil)
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
		assert.Equal(t, "category not found", responseBody.Data)
	})
}

func TestUpdateCategory(t *testing.T) {
	db := ConnectDBTest()
	DeleteDBTest(db)
	router := SetupRouterTest(db)
	defer db.Close()

	t.Run("success update category", func(t *testing.T) {
		ctx := context.Background()
		tx, err := db.Begin()
		helpers.PanicError(err)

		categoryRepository := repositoriesCategory.NewCategoryRepository()
		category := categoryRepository.Save(ctx, tx, domain.Category{Name: "category-5"})

		tx.Commit()

		categoryBody := strings.NewReader(`{
			"name": "category-6"
		}`)

		request := httptest.NewRequest(http.MethodPut, "http://localhost:8080/api/category/"+strconv.Itoa(category.Id), categoryBody)
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
		assert.Equal(t, "category-6", responseBody.Data.(map[string]interface{})["name"])
	})

	t.Run("not found update category", func(t *testing.T) {

		categoryBody := strings.NewReader(`{
			"name": "category-6"
		}`)

		request := httptest.NewRequest(http.MethodPut, "http://localhost:8080/api/category/1", categoryBody)
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
		assert.Equal(t, "category not found", responseBody.Data)
	})

	t.Run("bad request update category", func(t *testing.T) {
		ctx := context.Background()
		tx, err := db.Begin()
		helpers.PanicError(err)

		categoryRepository := repositoriesCategory.NewCategoryRepository()
		category := categoryRepository.Save(ctx, tx, domain.Category{Name: "category-5"})

		tx.Commit()

		categoryBody := strings.NewReader(`{
			"name": ""
		}`)

		request := httptest.NewRequest(http.MethodPut, "http://localhost:8080/api/category/"+strconv.Itoa(category.Id), categoryBody)
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
		assert.Equal(t, "Key: 'CategoryUpdateRequest.Name' Error:Field validation for 'Name' failed on the 'required' tag", responseBody.Data)
	})

}

func TestDeleteCategory(t *testing.T) {
	db := ConnectDBTest()
	DeleteDBTest(db)
	router := SetupRouterTest(db)
	defer db.Close()

	t.Run("success delete category", func(t *testing.T) {
		ctx := context.Background()
		tx, err := db.Begin()
		helpers.PanicError(err)

		categoryRepository := repositoriesCategory.NewCategoryRepository()
		category := categoryRepository.Save(ctx, tx, domain.Category{Name: "category-3"})

		tx.Commit()

		request := httptest.NewRequest(http.MethodDelete, "http://localhost:8080/api/category/"+strconv.Itoa(category.Id), nil)
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

	t.Run("not found delete category", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodDelete, "http://localhost:8080/api/category/1", nil)
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
