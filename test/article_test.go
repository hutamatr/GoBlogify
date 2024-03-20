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

	"github.com/hutamatr/go-blog-api/helpers"
	"github.com/hutamatr/go-blog-api/model/domain"
	"github.com/hutamatr/go-blog-api/model/web"

	repositoriesA "github.com/hutamatr/go-blog-api/repositories/article"
	repositoriesC "github.com/hutamatr/go-blog-api/repositories/category"
	"github.com/stretchr/testify/assert"
)

func createCategoryTestArticle(db *sql.DB) domain.Category {
	ctx := context.Background()
	tx, err := db.Begin()
	helpers.PanicError(err)
	defer tx.Commit()

	categoryRepository := repositoriesC.NewCategoryRepository()
	category := categoryRepository.Save(ctx, tx, domain.Category{Name: "category-3"})

	return category
}

func TestCreateArticle(t *testing.T) {
	db := ConnectDBTest()
	DeleteDBTest(db)
	router := SetupRouterTest(db)
	defer db.Close()

	category := createCategoryTestArticle(db)

	t.Run("success create article", func(t *testing.T) {
		articleBody := strings.NewReader(`{
			"title": "article-1",
			"body": "body-1",
			"author": "author-1",
			"published": true,
			"category_id": ` + strconv.Itoa(category.Id) + `
		}`)

		request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/article", articleBody)
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
		assert.Equal(t, "article-1", responseBody.Data.(map[string]interface{})["title"])
		assert.Equal(t, "body-1", responseBody.Data.(map[string]interface{})["body"])
		assert.Equal(t, "author-1", responseBody.Data.(map[string]interface{})["author"])
		assert.Equal(t, true, responseBody.Data.(map[string]interface{})["published"])
	})

	t.Run("bad request create article", func(t *testing.T) {
		articleBody := strings.NewReader(`{
			"title": "",
			"body": "body-1",
			"author": "author-1",
			"published": true,
			"category_id": ` + strconv.Itoa(category.Id) + `
		}`)

		request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/article", articleBody)
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

func TestFindAllArticle(t *testing.T) {
	db := ConnectDBTest()
	DeleteDBTest(db)
	router := SetupRouterTest(db)
	defer db.Close()

	category := createCategoryTestArticle(db)

	t.Run("success find all article", func(t *testing.T) {
		ctx := context.Background()
		tx, err := db.Begin()
		helpers.PanicError(err)

		articleRepository := repositoriesA.NewArticleRepository()
		article := articleRepository.Save(ctx, tx, domain.Article{
			Title:       "Article-3",
			Body:        "Body-3",
			Author:      "Author-3",
			Published:   true,
			Category_Id: category.Id,
		})

		tx.Commit()

		request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/article", nil)
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
		assert.Equal(t, article.Title, responseBody.Data.([]interface{})[0].(map[string]interface{})["title"])

	})

	t.Run("empty find all article", func(t *testing.T) {
		DeleteDBTest(db)

		request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/article", nil)
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

func TestFindByIdArticle(t *testing.T) {
	db := ConnectDBTest()
	DeleteDBTest(db)
	router := SetupRouterTest(db)
	defer db.Close()

	category := createCategoryTestArticle(db)

	t.Run("success find by id article", func(t *testing.T) {
		ctx := context.Background()
		tx, err := db.Begin()
		helpers.PanicError(err)

		articleRepository := repositoriesA.NewArticleRepository()
		article := articleRepository.Save(ctx, tx, domain.Article{
			Title:       "Article-4",
			Body:        "Body-4",
			Author:      "Author-4",
			Published:   true,
			Category_Id: category.Id,
		})

		tx.Commit()

		request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/article/"+strconv.Itoa(article.Id), nil)
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
		assert.Equal(t, article.Title, responseBody.Data.(map[string]interface{})["title"])
	})

	t.Run("not found find by id article", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/article/1", nil)
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
		assert.Equal(t, "article not found", responseBody.Data)
	})
}

func TestUpdateArticle(t *testing.T) {
	db := ConnectDBTest()
	DeleteDBTest(db)
	router := SetupRouterTest(db)
	defer db.Close()

	category := createCategoryTestArticle(db)

	t.Run("success update article", func(t *testing.T) {
		ctx := context.Background()
		tx, err := db.Begin()
		helpers.PanicError(err)

		articleRepository := repositoriesA.NewArticleRepository()
		article := articleRepository.Save(ctx, tx, domain.Article{
			Title:       "Article-5",
			Body:        "Body-5",
			Author:      "Author-5",
			Published:   true,
			Category_Id: category.Id,
		})

		tx.Commit()

		articleBody := strings.NewReader(`{
			"title" : "article-1",
			"body" : "body-1",
			"author" : "author-1",
			"published" : true,
			"category_id" : ` + strconv.Itoa(category.Id) + `
		}`)

		request := httptest.NewRequest(http.MethodPut, "http://localhost:8080/api/article/"+strconv.Itoa(article.Id), articleBody)
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
		assert.Equal(t, "article-1", responseBody.Data.(map[string]interface{})["title"])
		assert.Equal(t, "body-1", responseBody.Data.(map[string]interface{})["body"])

	})

	t.Run("not found update article", func(t *testing.T) {

		articleBody := strings.NewReader(`{
			"title" : "article-1",
			"body" : "body-1",
			"author" : "author-1",
			"published" : true,
			"category_id" : ` + strconv.Itoa(category.Id) + `
		}`)

		request := httptest.NewRequest(http.MethodPut, "http://localhost:8080/api/article/1", articleBody)
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
		assert.Equal(t, "article not found", responseBody.Data)
	})

	t.Run("bad request update article", func(t *testing.T) {
		ctx := context.Background()
		tx, err := db.Begin()
		helpers.PanicError(err)

		articleRepository := repositoriesA.NewArticleRepository()
		article := articleRepository.Save(ctx, tx, domain.Article{
			Title:       "Article-5",
			Body:        "Body-5",
			Author:      "Author-5",
			Published:   true,
			Category_Id: category.Id,
		})

		tx.Commit()

		articleBody := strings.NewReader(`{
			"title" : "",
			"body" : "body-1",
			"author" : "author-1",
			"published" : true,
			"category_id" : ` + strconv.Itoa(category.Id) + `
		}`)

		request := httptest.NewRequest(http.MethodPut, "http://localhost:8080/api/article/"+strconv.Itoa(article.Id), articleBody)
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

func TestDeleteArticle(t *testing.T) {
	db := ConnectDBTest()
	DeleteDBTest(db)
	router := SetupRouterTest(db)
	defer db.Close()

	category := createCategoryTestArticle(db)

	t.Run("success delete article", func(t *testing.T) {
		ctx := context.Background()
		tx, err := db.Begin()
		helpers.PanicError(err)

		articleRepository := repositoriesA.NewArticleRepository()
		article := articleRepository.Save(ctx, tx, domain.Article{
			Title:       "Article-5",
			Body:        "Body-5",
			Author:      "Author-5",
			Published:   true,
			Category_Id: category.Id,
		})

		tx.Commit()

		request := httptest.NewRequest(http.MethodDelete, "http://localhost:8080/api/article/"+strconv.Itoa(article.Id), nil)
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

	t.Run("not found delete article", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodDelete, "http://localhost:8080/api/article/100", nil)
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
