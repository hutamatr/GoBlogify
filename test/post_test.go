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

	"github.com/hutamatr/GoBlogify/category"
	"github.com/hutamatr/GoBlogify/helpers"
	"github.com/hutamatr/GoBlogify/post"
	"github.com/stretchr/testify/assert"
)

func createCategoryTestPost(db *sql.DB) category.Category {
	ctx := context.Background()
	tx, err := db.Begin()
	helpers.PanicError(err, "failed to begin transaction")
	defer tx.Commit()

	categoryRepository := category.NewCategoryRepository()
	category := categoryRepository.Save(ctx, tx, category.Category{Name: "category-3"})

	return category
}

func TestCreatePost(t *testing.T) {
	db := ConnectDBTest()
	DeleteDBTest(db)
	router := SetupRouterTest(db)
	defer db.Close()

	category := createCategoryTestPost(db)
	user, accessToken := createUserTestUser(db)

	t.Run("success create post", func(t *testing.T) {
		postBody := strings.NewReader(`{
			"title": "post-1",
			"body": "body-1",
			"user_id": ` + strconv.Itoa(user.Id) + `,
			"published": true,
			"category_id": ` + strconv.Itoa(category.Id) + `
		}`)

		request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/post", postBody)
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
		assert.Equal(t, "post-1", responseBody.Data.(map[string]interface{})["title"])
		assert.Equal(t, "body-1", responseBody.Data.(map[string]interface{})["body"])
		assert.Equal(t, true, responseBody.Data.(map[string]interface{})["published"])
	})

	t.Run("bad request create post", func(t *testing.T) {
		postBody := strings.NewReader(`{
			"title": "",
			"body": "body-1",
			"user_id": ` + strconv.Itoa(user.Id) + `,
			"published": true,
			"category_id": ` + strconv.Itoa(category.Id) + `
		}`)

		request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/post", postBody)
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

func TestFindAllPost(t *testing.T) {
	db := ConnectDBTest()
	DeleteDBTest(db)
	router := SetupRouterTest(db)
	defer db.Close()

	category := createCategoryTestPost(db)
	user, accessToken := createUserTestUser(db)

	t.Run("success find all post", func(t *testing.T) {
		ctx := context.Background()
		tx, err := db.Begin()
		helpers.PanicError(err, "failed to begin transaction")

		postRepository := post.NewPostRepository()
		post := postRepository.Save(ctx, tx, post.Post{
			Title:       "Post-3",
			Body:        "Body-3",
			User_Id:     user.Id,
			Published:   true,
			Category_Id: category.Id,
		})

		tx.Commit()

		request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/post", nil)
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
		assert.Equal(t, post.Title, responseBody.Data.(map[string]interface{})["posts"].([]interface{})[0].(map[string]interface{})["title"])
	})

	t.Run("empty find all post", func(t *testing.T) {
		DeleteDBTest(db)

		request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/post", nil)
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
		assert.Nil(t, responseBody.Data.(map[string]interface{})["posts"])
	})
}

func TestFindByIdPost(t *testing.T) {
	db := ConnectDBTest()
	DeleteDBTest(db)
	router := SetupRouterTest(db)
	defer db.Close()

	category := createCategoryTestPost(db)
	user, accessToken := createUserTestUser(db)

	t.Run("success find by id post", func(t *testing.T) {
		ctx := context.Background()
		tx, err := db.Begin()
		helpers.PanicError(err, "failed to begin transaction")

		postRepository := post.NewPostRepository()
		post := postRepository.Save(ctx, tx, post.Post{
			Title:       "Post-4",
			Body:        "Body-4",
			User_Id:     user.Id,
			Published:   true,
			Category_Id: category.Id,
		})

		tx.Commit()

		request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/post/"+strconv.Itoa(post.Id), nil)
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
		assert.Equal(t, post.Title, responseBody.Data.(map[string]interface{})["title"])
	})

	t.Run("not found find by id post", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/post/1", nil)
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
		assert.Equal(t, "post not found", responseBody.Data)
	})
}

func TestUpdatePost(t *testing.T) {
	db := ConnectDBTest()
	DeleteDBTest(db)
	router := SetupRouterTest(db)
	defer db.Close()

	category := createCategoryTestPost(db)
	user, accessToken := createUserTestUser(db)

	t.Run("success update post", func(t *testing.T) {
		ctx := context.Background()
		tx, err := db.Begin()
		helpers.PanicError(err, "failed to begin transaction")

		postRepository := post.NewPostRepository()
		post := postRepository.Save(ctx, tx, post.Post{
			Title:       "Post-5",
			Body:        "Body-5",
			User_Id:     user.Id,
			Published:   true,
			Category_Id: category.Id,
		})

		tx.Commit()

		postBody := strings.NewReader(`{
			"title" : "post-1",
			"body" : "body-1",
			"user_id" : ` + strconv.Itoa(user.Id) + `,
			"published" : true,
			"category_id" : ` + strconv.Itoa(category.Id) + `
		}`)

		request := httptest.NewRequest(http.MethodPut, "http://localhost:8080/api/post/"+strconv.Itoa(post.Id), postBody)
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
		assert.Equal(t, "post-1", responseBody.Data.(map[string]interface{})["title"])
		assert.Equal(t, "body-1", responseBody.Data.(map[string]interface{})["body"])

	})

	t.Run("not found update post", func(t *testing.T) {

		postBody := strings.NewReader(`{
			"title" : "post-1",
			"body" : "body-1",
			"user_id" : ` + strconv.Itoa(user.Id) + `,
			"published" : true,
			"category_id" : ` + strconv.Itoa(category.Id) + `
		}`)

		request := httptest.NewRequest(http.MethodPut, "http://localhost:8080/api/post/1", postBody)
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
		assert.Equal(t, "post not found", responseBody.Data)
	})

	t.Run("bad request update post", func(t *testing.T) {
		ctx := context.Background()
		tx, err := db.Begin()
		helpers.PanicError(err, "failed to begin transaction")

		postRepository := post.NewPostRepository()
		post := postRepository.Save(ctx, tx, post.Post{
			Title:       "Post-5",
			Body:        "Body-5",
			User_Id:     user.Id,
			Published:   true,
			Category_Id: category.Id,
		})

		tx.Commit()

		postBody := strings.NewReader(`{
			"title" : "",
			"body" : "body-1",
			"author" : "author-1",
			"published" : true,
			"category_id" : ` + strconv.Itoa(category.Id) + `
		}`)

		request := httptest.NewRequest(http.MethodPut, "http://localhost:8080/api/post/"+strconv.Itoa(post.Id), postBody)
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

func TestDeletePost(t *testing.T) {
	db := ConnectDBTest()
	DeleteDBTest(db)
	router := SetupRouterTest(db)
	defer db.Close()

	category := createCategoryTestPost(db)
	user, accessToken := createUserTestUser(db)

	t.Run("success delete post", func(t *testing.T) {
		ctx := context.Background()
		tx, err := db.Begin()
		helpers.PanicError(err, "failed to begin transaction")

		postRepository := post.NewPostRepository()
		post := postRepository.Save(ctx, tx, post.Post{
			Title:       "Post-5",
			Body:        "Body-5",
			User_Id:     user.Id,
			Published:   true,
			Category_Id: category.Id,
		})

		tx.Commit()

		request := httptest.NewRequest(http.MethodDelete, "http://localhost:8080/api/post/"+strconv.Itoa(post.Id), nil)
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

	t.Run("not found delete post", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodDelete, "http://localhost:8080/api/post/100", nil)
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
