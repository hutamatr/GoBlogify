package post_image

import "time"

type PostImageResponse struct {
	Id        int       `json:"id"`
	PostId    int       `json:"post_id"`
	Images    []string  `json:"images"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToPostImageResponse(postImage PostImage) PostImageResponse {
	return PostImageResponse{
		Id:        postImage.Id,
		PostId:    postImage.Post_Id,
		Images:    []string{postImage.Image_1, postImage.Image_2, postImage.Image_3},
		CreatedAt: postImage.Created_At,
		UpdatedAt: postImage.Updated_At,
	}
}
