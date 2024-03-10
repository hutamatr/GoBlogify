package go_blog

type ArticleCreateRequest struct {
	Title       string `json:"title" validate:"required,min=1,max=255"`
	Body        string `json:"body" validate:"required,min=1,max=1000"`
	Author      string `json:"author" validate:"required,min=1,max=100"`
	Published   bool   `json:"published" validate:"required"`
	Category_Id int    `json:"category_id"`
}

type ArticleUpdateRequest struct {
	Id          int    `json:"id" validate:"required"`
	Title       string `json:"title" validate:"required,min=1,max=255"`
	Body        string `json:"body" validate:"required,min=1,max=1000"`
	Author      string `json:"author" validate:"required,min=1,max=100"`
	Published   bool   `json:"published" validate:"required"`
	Deleted     bool   `json:"deleted"`
	Category_Id int    `json:"category_id"`
}
