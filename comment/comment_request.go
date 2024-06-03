package comment

type CommentCreateRequest struct {
	Comment string `json:"comment" validate:"required,min=1,max=500"`
	Post_Id int    `json:"post_id" validate:"required"`
	User_Id int    `json:"user_id" validate:"required"`
}

type CommentUpdateRequest struct {
	Id      int    `json:"id" validate:"required"`
	Comment string `json:"comment" validate:"required,min=1,max=500"`
}
