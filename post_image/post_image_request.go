package post_image

import "mime/multipart"

type PostImageRequest struct {
	Image_1      multipart.File `json:"image_1" validate:"required"`
	Image_Name_1 string         `json:"image_name_1" validate:"required"`
	Image_2      multipart.File `json:"image_2" validate:"required"`
	Image_Name_2 string         `json:"image_name_2" validate:"required"`
	Image_3      multipart.File `json:"image_3" validate:"required"`
	Image_Name_3 string         `json:"image_name_3" validate:"required"`
}

type PostImageUpdateRequest struct {
	Id           int            `json:"id"`
	Post_Id      int            `json:"post_id"`
	Image_1      multipart.File `json:"image_1"`
	Image_Name_1 string         `json:"image_name_1"`
	Image_2      multipart.File `json:"image_2"`
	Image_Name_2 string         `json:"image_name_2"`
	Image_3      multipart.File `json:"image_3"`
	Image_Name_3 string         `json:"image_name_3"`
}
