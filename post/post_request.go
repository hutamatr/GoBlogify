package post

import "mime/multipart"

type PostCreateRequest struct {
	Title        string         `json:"title" validate:"required,min=1,max=255"`
	Post_Body    string         `json:"post_body" validate:"required,min=1,max=1000"`
	Published    bool           `json:"published" validate:"required"`
	User_Id      int            `json:"user_id" validate:"required"`
	Category_Id  int            `json:"category_id" validate:"required"`
	Image_1      multipart.File `json:"image_1" validate:"required"`
	Image_Name_1 string         `json:"image_name_1" validate:"required"`
	Image_2      multipart.File `json:"image_2" validate:"required"`
	Image_Name_2 string         `json:"image_name_2" validate:"required"`
	Image_3      multipart.File `json:"image_3" validate:"required"`
	Image_Name_3 string         `json:"image_name_3" validate:"required"`
}

type PostUpdateRequest struct {
	Id           int            `json:"id" validate:"required"`
	User_Id      int            `json:"user_id" validate:"required"`
	Category_Id  int            `json:"category_id"`
	Title        string         `json:"title" validate:"required,min=1,max=255"`
	Post_Body    string         `json:"post_body" validate:"required,min=1,max=1000"`
	Published    bool           `json:"published" validate:"required"`
	Deleted      bool           `json:"deleted"`
	Image_1      multipart.File `json:"image_1"`
	Image_Name_1 string         `json:"image_name_1"`
	Image_2      multipart.File `json:"image_2"`
	Image_Name_2 string         `json:"image_name_2"`
	Image_3      multipart.File `json:"image_3"`
	Image_Name_3 string         `json:"image_name_3"`
}
