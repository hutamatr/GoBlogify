package test

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/hutamatr/GoBlogify/follow"
	"github.com/hutamatr/GoBlogify/helpers"
	"github.com/hutamatr/GoBlogify/role"
	"github.com/hutamatr/GoBlogify/user"
	"github.com/stretchr/testify/assert"
)

func createFollowTest(db *sql.DB, userId, toUserId int) follow.Follow {
	ctx := context.Background()
	tx, err := db.Begin()
	helpers.PanicError(err, "failed to begin transaction")
	defer tx.Commit()

	followRepository := follow.NewFollowRepository()
	follow := followRepository.Save(ctx, tx, follow.Follow{
		Follower_Id: userId,
		Followed_Id: toUserId,
	})

	return follow
}

func TestFollowUser(t *testing.T) {
	db := ConnectDBTest()
	DeleteDBTest(db)
	router := SetupRouterTest(db)
	defer db.Close()

	newUser1, accessToken := createUserTestUser(db)

	t.Run("success following user", func(t *testing.T) {
		ctx := context.Background()

		userRepository := user.NewUserRepository()
		roleRepository := role.NewRoleRepository()
		userService := user.NewUserService(userRepository, roleRepository, db, helpers.Validate)
		newUser2, _, _ := userService.SignUp(ctx, user.UserCreateRequest{Username: "userTest2", Email: "testing2@example.com", Password: "Password123!", Confirm_Password: "Password123!"})

		request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/users/"+strconv.Itoa(newUser1.Id)+"/follow/"+strconv.Itoa(newUser2.Id), nil)
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
	})

	t.Run("not found following user", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/users/"+strconv.Itoa(newUser1.Id)+"/follow/", nil)
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

func TestUnfollowUser(t *testing.T) {
	db := ConnectDBTest()
	DeleteDBTest(db)
	router := SetupRouterTest(db)
	defer db.Close()

	newUser1, accessToken := createUserTestUser(db)

	t.Run("success unfollow user", func(t *testing.T) {
		ctx := context.Background()

		userRepository := user.NewUserRepository()
		roleRepository := role.NewRoleRepository()
		userService := user.NewUserService(userRepository, roleRepository, db, helpers.Validate)
		newUser2, _, _ := userService.SignUp(ctx, user.UserCreateRequest{Username: "userTest2", Email: "testing2@example.com", Password: "Password123!", Confirm_Password: "Password123!"})

		followRepository := follow.NewFollowRepository()
		followService := follow.NewFollowService(followRepository, db)
		followedUser := followService.Following(ctx, newUser1.Id, newUser2.Id)

		request := httptest.NewRequest(http.MethodDelete, "http://localhost:8080/api/v1/users/"+strconv.Itoa(followedUser.Follower_Id)+"/unfollow/"+strconv.Itoa(followedUser.Followed_Id), nil)
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

	t.Run("not found unfollow user", func(t *testing.T) {
		DeleteDBTest(db)

		request := httptest.NewRequest(http.MethodDelete, "http://localhost:8080/api/v1/users/"+strconv.Itoa(newUser1.Id)+"/unfollow/", nil)
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

func TestFindAllFollowerByUser(t *testing.T) {
	db := ConnectDBTest()
	DeleteDBTest(db)
	router := SetupRouterTest(db)
	defer db.Close()

	newUser1, accessToken := createUserTestUser(db)

	t.Run("success find all follower", func(t *testing.T) {
		ctx := context.Background()

		userRepository := user.NewUserRepository()
		roleRepository := role.NewRoleRepository()
		userService := user.NewUserService(userRepository, roleRepository, db, helpers.Validate)
		newUser2, _, _ := userService.SignUp(ctx, user.UserCreateRequest{Username: "userTest2", Email: "testing2@example.com", Password: "Password123!", Confirm_Password: "Password123!"})

		followRepository := follow.NewFollowRepository()
		followService := follow.NewFollowService(followRepository, db)
		followService.Following(ctx, newUser1.Id, newUser2.Id)

		request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/v1/users/"+strconv.Itoa(newUser2.Id)+"/follower", nil)
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

		followerId := responseBody.Data.(map[string]interface{})["follower"].([]interface{})[0].(map[string]interface{})["follower_id"].(float64)

		assert.Equal(t, http.StatusOK, responseBody.Code)
		assert.Equal(t, "OK", responseBody.Status)
		assert.Equal(t, newUser1.Id, int(followerId))
	})

	t.Run("not found find all follower", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/v1/users/"+strconv.Itoa(newUser1.Id)+"/follower", nil)
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

func TestFindAllFollowedByUser(t *testing.T) {
	db := ConnectDBTest()
	DeleteDBTest(db)
	router := SetupRouterTest(db)
	defer db.Close()

	newUser1, accessToken := createUserTestUser(db)

	t.Run("success find all followed", func(t *testing.T) {
		ctx := context.Background()

		userRepository := user.NewUserRepository()
		roleRepository := role.NewRoleRepository()
		userService := user.NewUserService(userRepository, roleRepository, db, helpers.Validate)
		newUser2, _, _ := userService.SignUp(ctx, user.UserCreateRequest{Username: "userTest2", Email: "testing2@example.com", Password: "Password123!", Confirm_Password: "Password123!"})

		followRepository := follow.NewFollowRepository()
		followService := follow.NewFollowService(followRepository, db)
		followService.Following(ctx, newUser1.Id, newUser2.Id)

		request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/v1/users/"+strconv.Itoa(newUser1.Id)+"/following", nil)
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

		followedId := responseBody.Data.(map[string]interface{})["followed"].([]interface{})[0].(map[string]interface{})["followed_id"].(float64)

		assert.Equal(t, http.StatusOK, responseBody.Code)
		assert.Equal(t, "OK", responseBody.Status)
		assert.Equal(t, newUser2.Id, int(followedId))
	})

	t.Run("not found find all followed", func(t *testing.T) {
		ctx := context.Background()

		userRepository := user.NewUserRepository()
		roleRepository := role.NewRoleRepository()
		userService := user.NewUserService(userRepository, roleRepository, db, helpers.Validate)
		newUser3, _, _ := userService.SignUp(ctx, user.UserCreateRequest{Username: "userTest3", Email: "testing3@example.com", Password: "Password123!", Confirm_Password: "Password123!"})

		request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/v1/users/"+strconv.Itoa(newUser3.Id)+"/following", nil)
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
