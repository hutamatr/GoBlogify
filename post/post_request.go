package post

type PostCreateRequest struct {
	Title       string `json:"title" validate:"required,min=1,max=255"`
	Body        string `json:"body" validate:"required,min=1,max=1000"`
	Published   bool   `json:"published" validate:"required"`
	User_Id     int    `json:"user_id" validate:"required"`
	Category_Id int    `json:"category_id" validate:"required"`
}

type PostUpdateRequest struct {
	Id          int    `json:"id" validate:"required"`
	User_Id     int    `json:"user_id" validate:"required"`
	Category_Id int    `json:"category_id"`
	Title       string `json:"title" validate:"required,min=1,max=255"`
	Body        string `json:"body" validate:"required,min=1,max=1000"`
	Published   bool   `json:"published" validate:"required"`
	Deleted     bool   `json:"deleted"`
}
