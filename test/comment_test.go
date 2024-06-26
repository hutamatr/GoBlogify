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

	"github.com/hutamatr/GoBlogify/comment"
	"github.com/hutamatr/GoBlogify/helpers"
	"github.com/hutamatr/GoBlogify/post"
	"github.com/stretchr/testify/assert"
)

func createPostTestComment(db *sql.DB, userId int, categoryId int) post.PostJoin {
	ctx := context.Background()
	tx, err := db.Begin()
	helpers.PanicError(err, "failed to begin transaction")
	defer tx.Commit()

	postRepository := post.NewPostRepository()
	post := postRepository.Save(ctx, tx, post.Post{
		Title:       "Post-1",
		Body:        "Body-1",
		Published:   true,
		User_Id:     userId,
		Category_Id: categoryId,
	})

	return post
}

func TestCreateComment(t *testing.T) {
	db := ConnectDBTest()
	DeleteDBTest(db)
	router := SetupRouterTest(db)
	defer db.Close()

	user, accessToken := createUserTestUser(db)
	category := createCategoryTestPost(db)
	post := createPostTestComment(db, user.Id, category.Id)

	t.Run("success create comment", func(t *testing.T) {
		commentBody := strings.NewReader(`{
			"content": "comment-1",
			"user_id": ` + strconv.Itoa(user.Id) + `,
			"post_id": ` + strconv.Itoa(post.Id) + `
		}`)

		request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/comments", commentBody)
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
		assert.Equal(t, "comment-1", responseBody.Data.(map[string]interface{})["content"])
		assert.Equal(t, post.Id, int(responseBody.Data.(map[string]interface{})["post_id"].(float64)))
		assert.Equal(t, user.Id, int(responseBody.Data.(map[string]interface{})["user_id"].(float64)))
	})

	t.Run("bad request create comment", func(t *testing.T) {
		commentBody := strings.NewReader(`{
			"content": "",
			"user_id": ` + strconv.Itoa(user.Id) + `,
			"post_id": ` + strconv.Itoa(post.Id) + `
		}`)

		request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/comments", commentBody)
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", "Bearer "+accessToken)

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

func TestFindCommentByPost(t *testing.T) {
	db := ConnectDBTest()
	DeleteDBTest(db)
	router := SetupRouterTest(db)
	defer db.Close()

	category := createCategoryTestPost(db)
	user, accessToken := createUserTestUser(db)
	post := createPostTestComment(db, user.Id, category.Id)

	t.Run("success find comment by post", func(t *testing.T) {
		ctx := context.Background()
		tx, err := db.Begin()
		helpers.PanicError(err, "failed to begin transaction")

		commentRepository := comment.NewCommentRepository()
		comment := commentRepository.Save(ctx, tx, comment.Comment{
			Content: "comment-1",
			User_Id: user.Id,
			Post_Id: post.Id,
		})

		tx.Commit()

		request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/v1/comments?postId="+strconv.Itoa(post.Id), nil)
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
		assert.Equal(t, comment.Content, responseBody.Data.(map[string]interface{})["comments"].([]interface{})[0].(map[string]interface{})["content"])

	})

	t.Run("empty find comment by post", func(t *testing.T) {
		DeleteDBTest(db)

		request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/v1/comments?postId="+strconv.Itoa(post.Id), nil)
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", "Bearer "+accessToken)

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, request)

		response := recorder.Result()

		assert.Equal(t, http.StatusNotFound, response.StatusCode)

		body, err := io.ReadAll(response.Body)

		var responseBody helpers.ErrorResponseJSON

		json.Unmarshal(body, &responseBody)

		helpers.PanicError(err, "failed to read response body")

		assert.Equal(t, http.StatusNotFound, responseBody.Code)
		assert.Equal(t, "NOT FOUND", responseBody.Status)
	})
}

func TestFindCommentById(t *testing.T) {
	db := ConnectDBTest()
	DeleteDBTest(db)
	router := SetupRouterTest(db)
	defer db.Close()

	category := createCategoryTestPost(db)
	user, accessToken := createUserTestUser(db)
	post := createPostTestComment(db, user.Id, category.Id)

	t.Run("success find comment by id", func(t *testing.T) {
		ctx := context.Background()
		tx, err := db.Begin()
		helpers.PanicError(err, "failed to begin transaction")

		commentRepository := comment.NewCommentRepository()
		comment := commentRepository.Save(ctx, tx, comment.Comment{
			Content: "comment-1",
			User_Id: user.Id,
			Post_Id: post.Id,
		})

		tx.Commit()

		request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/v1/comments/"+strconv.Itoa(comment.Id), nil)
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
		assert.Equal(t, comment.Content, responseBody.Data.(map[string]interface{})["content"])
	})

	t.Run("not found find comment by id", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/v1/comments/1", nil)
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", "Bearer "+accessToken)

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, request)

		response := recorder.Result()

		assert.Equal(t, http.StatusNotFound, response.StatusCode)

		body, err := io.ReadAll(response.Body)

		var responseBody helpers.ErrorResponseJSON

		json.Unmarshal(body, &responseBody)

		helpers.PanicError(err, "failed to read response body")

		assert.Equal(t, http.StatusNotFound, responseBody.Code)
		assert.Equal(t, "NOT FOUND", responseBody.Status)
		assert.Equal(t, "comment not found", responseBody.Error)
	})
}

func TestUpdateComment(t *testing.T) {
	db := ConnectDBTest()
	DeleteDBTest(db)
	router := SetupRouterTest(db)
	defer db.Close()

	category := createCategoryTestPost(db)
	user, accessToken := createUserTestUser(db)
	post := createPostTestComment(db, user.Id, category.Id)

	t.Run("success update comment", func(t *testing.T) {
		ctx := context.Background()
		tx, err := db.Begin()
		helpers.PanicError(err, "failed to begin transaction")

		commentRepository := comment.NewCommentRepository()
		comment := commentRepository.Save(ctx, tx, comment.Comment{
			Content: "comment-1",
			User_Id: user.Id,
			Post_Id: post.Id,
		})

		tx.Commit()

		commentBody := strings.NewReader(`{
			"content": "comment-2",
			"user_id": ` + strconv.Itoa(user.Id) + `,
			"post_id": ` + strconv.Itoa(post.Id) + `
		}`)

		request := httptest.NewRequest(http.MethodPut, "http://localhost:8080/api/v1/comments/"+strconv.Itoa(comment.Id), commentBody)
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
		assert.Equal(t, "comment-2", responseBody.Data.(map[string]interface{})["content"])
	})

	t.Run("not found update comment", func(t *testing.T) {
		commentBody := strings.NewReader(`{
			"content": "comment-2",
			"user_id": ` + strconv.Itoa(user.Id) + `,
			"post_id": ` + strconv.Itoa(post.Id) + `
		}`)

		request := httptest.NewRequest(http.MethodPut, "http://localhost:8080/api/v1/comments/1", commentBody)
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", "Bearer "+accessToken)

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, request)

		response := recorder.Result()

		assert.Equal(t, http.StatusNotFound, response.StatusCode)

		body, err := io.ReadAll(response.Body)

		var responseBody helpers.ErrorResponseJSON

		json.Unmarshal(body, &responseBody)

		helpers.PanicError(err, "failed to read response body")

		assert.Equal(t, http.StatusNotFound, responseBody.Code)
		assert.Equal(t, "NOT FOUND", responseBody.Status)
		assert.Equal(t, "comment not found", responseBody.Error)
	})

	t.Run("bad request update comment", func(t *testing.T) {
		ctx := context.Background()
		tx, err := db.Begin()
		helpers.PanicError(err, "failed to begin transaction")

		commentRepository := comment.NewCommentRepository()
		comment := commentRepository.Save(ctx, tx, comment.Comment{
			Content: "comment-1",
			User_Id: user.Id,
			Post_Id: post.Id,
		})

		tx.Commit()

		commentBody := strings.NewReader(`{
			"content": "",
			"user_id": ` + strconv.Itoa(user.Id) + `,
			"post_id": ` + strconv.Itoa(post.Id) + `
		}`)

		request := httptest.NewRequest(http.MethodPut, "http://localhost:8080/api/v1/comments/"+strconv.Itoa(comment.Id), commentBody)
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", "Bearer "+accessToken)

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

func TestDeleteComment(t *testing.T) {
	db := ConnectDBTest()
	DeleteDBTest(db)
	router := SetupRouterTest(db)
	defer db.Close()

	category := createCategoryTestPost(db)
	user, accessToken := createUserTestUser(db)
	post := createPostTestComment(db, user.Id, category.Id)

	t.Run("success delete comment", func(t *testing.T) {
		ctx := context.Background()
		tx, err := db.Begin()
		helpers.PanicError(err, "failed to begin transaction")

		commentRepository := comment.NewCommentRepository()
		comment := commentRepository.Save(ctx, tx, comment.Comment{
			Content: "comment-1",
			User_Id: user.Id,
			Post_Id: post.Id,
		})

		tx.Commit()

		request := httptest.NewRequest(http.MethodDelete, "http://localhost:8080/api/v1/comments/"+strconv.Itoa(comment.Id), nil)
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

	t.Run("not found delete comment", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodDelete, "http://localhost:8080/api/v1/comments/100", nil)
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", "Bearer "+accessToken)

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, request)

		response := recorder.Result()

		assert.Equal(t, http.StatusNotFound, response.StatusCode)

		body, err := io.ReadAll(response.Body)

		var responseBody helpers.ErrorResponseJSON

		json.Unmarshal(body, &responseBody)

		helpers.PanicError(err, "failed to read response body")

		assert.Equal(t, http.StatusNotFound, responseBody.Code)
		assert.Equal(t, "NOT FOUND", responseBody.Status)
	})
}
