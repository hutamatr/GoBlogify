package user

import (
	"time"
)

type UserResponse struct {
	Id         int       `json:"id"`
	Role_Id    int       `json:"role_id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	First_Name string    `json:"first_name"`
	Last_Name  string    `json:"last_name"`
	Created_At time.Time `json:"created_at"`
	Updated_At time.Time `json:"updated_at"`
	Deleted_At time.Time `json:"deleted_at"`
	Following  int       `json:"following"`
	Follower   int       `json:"follower"`
}

func ToUserResponse(user UserJoin) UserResponse {
	return UserResponse{
		Id:         user.Id,
		Role_Id:    user.Role_Id,
		Username:   user.Username,
		Email:      user.Email,
		First_Name: user.First_Name,
		Last_Name:  user.Last_Name,
		Created_At: user.Created_At,
		Updated_At: user.Updated_At,
		Deleted_At: user.Deleted_At,
		Following:  user.Following,
		Follower:   user.Follower,
	}
}

type UserCommentResponse struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func ToUserCommentResponse(user User) UserCommentResponse {
	return UserCommentResponse{
		Id:       user.Id,
		Username: user.Username,
		Email:    user.Email,
	}
}

type UserFollowResponse struct {
	Id         int    `json:"id"`
	Username   string `json:"username"`
	First_Name string `json:"first_name"`
	Last_Name  string `json:"last_name"`
}

func ToUserFollowResponse(user User) UserFollowResponse {
	return UserFollowResponse{
		Id:         user.Id,
		Username:   user.Username,
		First_Name: user.First_Name,
		Last_Name:  user.Last_Name,
	}
}
