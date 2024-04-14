package follow

import (
	"time"

	"github.com/hutamatr/GoBlogify/user"
)

type FollowResponse struct {
	Id          int       `json:"id"`
	Follower_Id int       `json:"follower_id"`
	Followed_Id int       `json:"followed_id"`
	Created_At  time.Time `json:"created_at"`
	Updated_At  time.Time `json:"updated_at"`
}

func ToFollowResponse(follow Follow) FollowResponse {
	return FollowResponse{
		Id:          follow.Id,
		Follower_Id: follow.Follower_Id,
		Followed_Id: follow.Followed_Id,
		Created_At:  follow.Created_At,
		Updated_At:  follow.Updated_At,
	}
}

type FollowJoinResponse struct {
	Id          int                     `json:"id"`
	Follower_Id int                     `json:"follower_id"`
	Followed_Id int                     `json:"followed_id"`
	Created_At  time.Time               `json:"created_at"`
	Updated_At  time.Time               `json:"updated_at"`
	User        user.UserFollowResponse `json:"user"`
}

func ToFollowJoinResponse(follow FollowJoin) FollowJoinResponse {
	return FollowJoinResponse{
		Id:          follow.Id,
		Follower_Id: follow.Follower_Id,
		Followed_Id: follow.Followed_Id,
		Created_At:  follow.Created_At,
		Updated_At:  follow.Updated_At,
		User:        user.ToUserFollowResponse(follow.User),
	}
}
