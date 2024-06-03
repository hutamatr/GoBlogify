package test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	_ "embed"

	"github.com/hutamatr/GoBlogify/category"
	"github.com/hutamatr/GoBlogify/helpers"
	"github.com/hutamatr/GoBlogify/post"
	"github.com/hutamatr/GoBlogify/post_image"
	"github.com/hutamatr/GoBlogify/role"
	"github.com/hutamatr/GoBlogify/user"
	"github.com/stretchr/testify/assert"
)

//go:embed test.webp
var testImage []byte

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
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)

		writer.WriteField("title", "post-1")
		writer.WriteField("post_body", "body-1")
		writer.WriteField("published", "true")
		writer.WriteField("user_id", strconv.Itoa(user.Id))
		writer.WriteField("category_id", strconv.Itoa(category.Id))

		AddFile(t, writer, "attachment", "test.webp")
		AddFile(t, writer, "attachment", "test.webp")
		AddFile(t, writer, "attachment", "test.webp")

		writer.Close()

		request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/posts", body)
		request.Header.Add("Content-Type", writer.FormDataContentType())
		request.Header.Add("Authorization", "Bearer "+accessToken)

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, request)

		response := recorder.Result()

		assert.Equal(t, http.StatusCreated, response.StatusCode)

		bodyRes, err := io.ReadAll(response.Body)

		var responseBody helpers.ResponseJSON

		json.Unmarshal(bodyRes, &responseBody)

		helpers.PanicError(err, "failed to read response body")

		assert.Equal(t, http.StatusCreated, responseBody.Code)
		assert.Equal(t, "CREATED", responseBody.Status)
		assert.Equal(t, "post-1", responseBody.Data.(map[string]interface{})["title"])
		assert.Equal(t, "body-1", responseBody.Data.(map[string]interface{})["post_body"])
		assert.Equal(t, true, responseBody.Data.(map[string]interface{})["published"])
	})

	t.Run("bad request create post", func(t *testing.T) {
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)

		writer.WriteField("title", "")
		writer.WriteField("post_body", "")
		writer.WriteField("published", "true")
		writer.WriteField("user_id", strconv.Itoa(user.Id))
		writer.WriteField("category_id", strconv.Itoa(category.Id))

		AddFile(t, writer, "attachment", "test.webp")
		AddFile(t, writer, "attachment", "test.webp")
		AddFile(t, writer, "attachment", "test.webp")

		writer.Close()

		request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/posts", body)
		request.Header.Add("Content-Type", writer.FormDataContentType())
		request.Header.Add("Authorization", "Bearer "+accessToken)

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, request)

		response := recorder.Result()

		assert.Equal(t, http.StatusBadRequest, response.StatusCode)

		bodyRes, err := io.ReadAll(response.Body)

		var responseBody helpers.ErrorResponseJSON

		json.Unmarshal(bodyRes, &responseBody)

		helpers.PanicError(err, "failed to read response body")

		assert.Equal(t, http.StatusBadRequest, responseBody.Code)
		assert.Equal(t, "BAD REQUEST", responseBody.Status)
	})
}

func TestFindAllPostByUser(t *testing.T) {
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

		files := []string{"test.webp", "test.webp", "test.webp"}
		var imageUrls []string

		for _, file := range files {
			imageUrl, err := helpers.UploadToCloudinary(file, "test.webp")
			assert.NoError(t, err)

			imageUrls = append(imageUrls, imageUrl)
		}

		postRepository := post.NewPostRepository()
		postCreated := postRepository.Save(ctx, tx, post.Post{
			Title:       "Post5",
			Post_Body:   "Body5",
			User_Id:     user.Id,
			Published:   true,
			Category_Id: category.Id,
		})

		postImageRepository := post_image.NewPostImageRepository()
		postImageRepository.Save(ctx, tx, post_image.PostImage{
			Post_Id:      postCreated.Id,
			Image_1:      imageUrls[0],
			Image_Name_1: "test.web",
			Image_2:      imageUrls[1],
			Image_Name_2: "test.webp",
			Image_3:      imageUrls[2],
			Image_Name_3: "test.webp",
		})

		tx.Commit()

		request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/v1/posts/"+strconv.Itoa(user.Id), nil)
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
		assert.Equal(t, postCreated.Title, responseBody.Data.(map[string]interface{})["posts"].([]interface{})[0].(map[string]interface{})["title"])
	})

	t.Run("empty find all post by user", func(t *testing.T) {
		DeleteDBTest(db)

		request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/v1/posts/"+strconv.Itoa(user.Id), nil)
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

func TestFindAllPostByFollowed(t *testing.T) {
	db := ConnectDBTest()
	DeleteDBTest(db)
	router := SetupRouterTest(db)
	defer db.Close()

	category := createCategoryTestPost(db)
	newUser1, accessToken1 := createUserTestUser(db)

	t.Run("success find all post by followed", func(t *testing.T) {
		ctx := context.Background()
		tx, err := db.Begin()
		helpers.PanicError(err, "failed to begin transaction")

		files := []string{"test.webp", "test.webp", "test.webp"}
		var imageUrls []string

		for _, file := range files {
			imageUrl, err := helpers.UploadToCloudinary(file, "test.webp")
			assert.NoError(t, err)

			imageUrls = append(imageUrls, imageUrl)
		}

		postRepository := post.NewPostRepository()
		postCreated := postRepository.Save(ctx, tx, post.Post{
			Title:       "Post5",
			Post_Body:   "Body5",
			User_Id:     newUser1.Id,
			Published:   true,
			Category_Id: category.Id,
		})

		postImageRepository := post_image.NewPostImageRepository()
		postImageRepository.Save(ctx, tx, post_image.PostImage{
			Post_Id:      postCreated.Id,
			Image_1:      imageUrls[0],
			Image_Name_1: "test.web",
			Image_2:      imageUrls[1],
			Image_Name_2: "test.webp",
			Image_3:      imageUrls[2],
			Image_Name_3: "test.webp",
		})
		tx.Commit()

		userRepository := user.NewUserRepository()
		roleRepository := role.NewRoleRepository()
		userService := user.NewUserService(userRepository, roleRepository, db, helpers.Validate)
		newUser2, accessToken2, _ := userService.SignUp(ctx, user.UserCreateRequest{Username: "userTest2", Email: "testing2@example.com", Password: "Password123!", Confirm_Password: "Password123!"})

		followUser := createFollowTest(db, newUser2.Id, newUser1.Id)

		request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/v1/posts/"+strconv.Itoa(followUser.Follower_Id)+"/following", nil)
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", "Bearer "+accessToken2)

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
		assert.Equal(t, postCreated.Title, responseBody.Data.(map[string]interface{})["posts"].([]interface{})[0].(map[string]interface{})["title"])
		assert.Equal(t, postCreated.Post_Body, responseBody.Data.(map[string]interface{})["posts"].([]interface{})[0].(map[string]interface{})["post_body"])
	})

	t.Run("not found find all post by followed", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/v1/posts/"+strconv.Itoa(newUser1.Id)+"/following", nil)
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", "Bearer "+accessToken1)

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

func TestFindPostById(t *testing.T) {
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

		files := []string{"test.webp", "test.webp", "test.webp"}
		var imageUrls []string

		for _, file := range files {
			imageUrl, err := helpers.UploadToCloudinary(file, "test.webp")
			assert.NoError(t, err)

			imageUrls = append(imageUrls, imageUrl)
		}

		postRepository := post.NewPostRepository()
		postCreated := postRepository.Save(ctx, tx, post.Post{
			Title:       "Post5",
			Post_Body:   "Body5",
			User_Id:     user.Id,
			Published:   true,
			Category_Id: category.Id,
		})

		postImageRepository := post_image.NewPostImageRepository()
		postImageRepository.Save(ctx, tx, post_image.PostImage{
			Post_Id:      postCreated.Id,
			Image_1:      imageUrls[0],
			Image_Name_1: "test.web",
			Image_2:      imageUrls[1],
			Image_Name_2: "test.webp",
			Image_3:      imageUrls[2],
			Image_Name_3: "test.webp",
		})

		tx.Commit()

		request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/v1/post/"+strconv.Itoa(postCreated.Id), nil)
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
		assert.Equal(t, postCreated.Title, responseBody.Data.(map[string]interface{})["title"])
	})

	t.Run("not found find post by id", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/v1/post/100", nil)
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
		assert.Equal(t, "post not found", responseBody.Error)
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

		files := []string{"test.webp", "test.webp", "test.webp"}
		var imageUrls []string

		for _, file := range files {
			imageUrl, err := helpers.UploadToCloudinary(file, "test.webp")
			assert.NoError(t, err)

			imageUrls = append(imageUrls, imageUrl)
		}

		postRepository := post.NewPostRepository()
		postCreated := postRepository.Save(ctx, tx, post.Post{
			Title:       "Post5",
			Post_Body:   "Body5",
			User_Id:     user.Id,
			Published:   true,
			Category_Id: category.Id,
		})

		postImageRepository := post_image.NewPostImageRepository()
		postImageRepository.Save(ctx, tx, post_image.PostImage{
			Post_Id:      postCreated.Id,
			Image_1:      imageUrls[0],
			Image_Name_1: "test.web",
			Image_2:      imageUrls[1],
			Image_Name_2: "test.webp",
			Image_3:      imageUrls[2],
			Image_Name_3: "test.webp",
		})

		tx.Commit()

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		// fileImg, err := writer.CreateFormFile("file", "test.webp")
		// assert.Nil(t, err)
		// fileImg.Write(testImage)

		writer.WriteField("title", "post1")
		writer.WriteField("post_body", "body1")
		writer.WriteField("published", "true")
		writer.WriteField("user_id", strconv.Itoa(user.Id))
		writer.WriteField("category_id", strconv.Itoa(category.Id))

		AddFile(t, writer, "attachment", "test.webp")
		AddFile(t, writer, "attachment", "test.webp")
		AddFile(t, writer, "attachment", "test.webp")

		writer.Close()

		request := httptest.NewRequest(http.MethodPut, "http://localhost:8080/api/v1/posts/"+strconv.Itoa(postCreated.Id), body)
		request.Header.Add("Content-Type", writer.FormDataContentType())
		request.Header.Add("Authorization", "Bearer "+accessToken)

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, request)

		response := recorder.Result()

		assert.Equal(t, http.StatusOK, response.StatusCode)

		bodyRes, err := io.ReadAll(response.Body)

		var responseBody helpers.ResponseJSON

		json.Unmarshal(bodyRes, &responseBody)

		helpers.PanicError(err, "failed to read response body")

		assert.Equal(t, http.StatusOK, responseBody.Code)
		assert.Equal(t, "UPDATED", responseBody.Status)
		assert.Equal(t, "post1", responseBody.Data.(map[string]interface{})["title"])
		assert.Equal(t, "body1", responseBody.Data.(map[string]interface{})["post_body"])

	})

	t.Run("not found update post", func(t *testing.T) {

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)

		writer.WriteField("title", "post-1")
		writer.WriteField("post_body", "body-1")
		writer.WriteField("published", "true")
		writer.WriteField("user_id", strconv.Itoa(user.Id))
		writer.WriteField("category_id", strconv.Itoa(category.Id))

		AddFile(t, writer, "attachment", "test.webp")
		AddFile(t, writer, "attachment", "test.webp")
		AddFile(t, writer, "attachment", "test.webp")

		writer.Close()

		request := httptest.NewRequest(http.MethodPut, "http://localhost:8080/api/v1/posts/1", body)
		request.Header.Add("Content-Type", writer.FormDataContentType())
		request.Header.Add("Authorization", "Bearer "+accessToken)

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, request)

		response := recorder.Result()

		assert.Equal(t, http.StatusNotFound, response.StatusCode)

		bodyRes, err := io.ReadAll(response.Body)

		var responseBody helpers.ErrorResponseJSON

		json.Unmarshal(bodyRes, &responseBody)

		helpers.PanicError(err, "failed to read response body")

		assert.Equal(t, http.StatusNotFound, responseBody.Code)
		assert.Equal(t, "NOT FOUND", responseBody.Status)
		assert.Equal(t, "post not found", responseBody.Error)
	})

	t.Run("bad request update post", func(t *testing.T) {
		ctx := context.Background()
		tx, err := db.Begin()
		helpers.PanicError(err, "failed to begin transaction")

		files := []string{"test.webp", "test.webp", "test.webp"}
		var imageUrls []string

		for _, file := range files {
			imageUrl, err := helpers.UploadToCloudinary(file, "test.webp")
			assert.NoError(t, err)

			imageUrls = append(imageUrls, imageUrl)
		}

		postRepository := post.NewPostRepository()
		postCreated := postRepository.Save(ctx, tx, post.Post{
			Title:       "Post5",
			Post_Body:   "Body5",
			User_Id:     user.Id,
			Published:   true,
			Category_Id: category.Id,
		})

		postImageRepository := post_image.NewPostImageRepository()
		postImageRepository.Save(ctx, tx, post_image.PostImage{
			Post_Id:      postCreated.Id,
			Image_1:      imageUrls[0],
			Image_Name_1: "test.web",
			Image_2:      imageUrls[1],
			Image_Name_2: "test.webp",
			Image_3:      imageUrls[2],
			Image_Name_3: "test.webp",
		})

		tx.Commit()

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		fileImg, err := writer.CreateFormFile("file", "test.webp")
		assert.Nil(t, err)
		fileImg.Write(testImage)

		writer.WriteField("title", "")
		writer.WriteField("post_body", "body1")
		writer.WriteField("published", "true")
		writer.WriteField("user_id", strconv.Itoa(user.Id))
		writer.WriteField("category_id", strconv.Itoa(category.Id))

		writer.Close()

		request := httptest.NewRequest(http.MethodPut, "http://localhost:8080/api/v1/posts/"+strconv.Itoa(postCreated.Id), body)
		request.Header.Add("Content-Type", writer.FormDataContentType())
		request.Header.Add("Authorization", "Bearer "+accessToken)

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, request)

		response := recorder.Result()

		assert.Equal(t, http.StatusBadRequest, response.StatusCode)

		bodyRes, err := io.ReadAll(response.Body)

		var responseBody helpers.ErrorResponseJSON

		json.Unmarshal(bodyRes, &responseBody)

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
			Post_Body:   "Body-5",
			User_Id:     user.Id,
			Published:   true,
			Category_Id: category.Id,
		})

		tx.Commit()

		request := httptest.NewRequest(http.MethodDelete, "http://localhost:8080/api/v1/posts/"+strconv.Itoa(post.Id), nil)
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
		request := httptest.NewRequest(http.MethodDelete, "http://localhost:8080/api/v1/posts/100", nil)
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
