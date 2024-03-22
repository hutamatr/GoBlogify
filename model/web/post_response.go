package web

import (
	"time"

	"github.com/hutamatr/GoBlogify/model/domain"
)

type PostResponse struct {
	Id         int              `json:"id"`
	Title      string           `json:"title"`
	Body       string           `json:"body"`
	Published  bool             `json:"published"`
	Deleted    bool             `json:"deleted"`
	Created_At time.Time        `json:"created_at"`
	Updated_At time.Time        `json:"updated_at"`
	Deleted_At time.Time        `json:"deleted_at"`
	User       UserResponse     `json:"user"`
	Category   CategoryResponse `json:"category"`
}

func ToPostResponse(post domain.PostJoin) PostResponse {
	return PostResponse{
		Id:         post.Id,
		Title:      post.Title,
		Body:       post.Body,
		Published:  post.Published,
		Deleted:    post.Deleted,
		Created_At: post.Created_At,
		Updated_At: post.Updated_At,
		Deleted_At: post.Deleted_At,
		User:       ToUserResponse(post.User),
		Category:   ToCategoryResponse(post.Category),
	}
}
