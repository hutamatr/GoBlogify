package post

import (
	"time"

	"github.com/hutamatr/GoBlogify/category"
	"github.com/hutamatr/GoBlogify/post_image"
	"github.com/hutamatr/GoBlogify/user"
)

type PostResponse struct {
	Id         int                          `json:"id"`
	Title      string                       `json:"title"`
	Post_Body  string                       `json:"post_body"`
	Published  bool                         `json:"published"`
	Deleted    bool                         `json:"deleted"`
	Created_At time.Time                    `json:"created_at"`
	Updated_At time.Time                    `json:"updated_at"`
	Deleted_At time.Time                    `json:"deleted_at"`
	User       user.UserResponse            `json:"user"`
	Category   category.CategoryResponse    `json:"category"`
	Images     post_image.PostImageResponse `json:"images"`
}

func ToPostResponse(post PostJoin, postImage post_image.PostImage) PostResponse {
	return PostResponse{
		Id:         post.Id,
		Title:      post.Title,
		Post_Body:  post.Post_Body,
		Published:  post.Published,
		Deleted:    post.Deleted,
		Created_At: post.Created_At,
		Updated_At: post.Updated_At,
		Deleted_At: post.Deleted_At,
		User:       user.ToUserResponse(post.User),
		Category:   category.ToCategoryResponse(post.Category),
		Images:     post_image.ToPostImageResponse(postImage),
	}
}

type PostResponseFollowed struct {
	Id         int                          `json:"id"`
	Title      string                       `json:"title"`
	Post_Body  string                       `json:"post_body"`
	Published  bool                         `json:"published"`
	Deleted    bool                         `json:"deleted"`
	Created_At time.Time                    `json:"created_at"`
	Updated_At time.Time                    `json:"updated_at"`
	Deleted_At time.Time                    `json:"deleted_at"`
	User       user.UserResponse            `json:"user"`
	Images     post_image.PostImageResponse `json:"images"`
}

func ToPostResponseFollowed(post PostJoinFollowed) PostResponseFollowed {
	return PostResponseFollowed{
		Id:         post.Id,
		Title:      post.Title,
		Post_Body:  post.Post_Body,
		Published:  post.Published,
		Deleted:    post.Deleted,
		Created_At: post.Created_At,
		Updated_At: post.Updated_At,
		Deleted_At: post.Deleted_At,
		User:       user.ToUserResponse(post.User),
		Images:     post_image.ToPostImageResponse(post.Images),
	}
}
