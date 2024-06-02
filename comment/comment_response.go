package comment

import (
	"time"

	"github.com/hutamatr/GoBlogify/user"
)

type CommentResponse struct {
	Id         int                      `json:"id"`
	Post_Id    int                      `json:"post_id"`
	User_Id    int                      `json:"user_id"`
	Comment    string                   `json:"comment"`
	Created_At time.Time                `json:"created_at"`
	Updated_At time.Time                `json:"updated_at"`
	User       user.UserCommentResponse `json:"user"`
}

func ToCommentResponse(comment CommentJoin) CommentResponse {
	return CommentResponse{
		Id:         comment.Id,
		Post_Id:    comment.Post_Id,
		User_Id:    comment.User_Id,
		Comment:    comment.Comment,
		Created_At: comment.Created_At,
		Updated_At: comment.Updated_At,
		User:       user.ToUserCommentResponse(comment.User),
	}
}
