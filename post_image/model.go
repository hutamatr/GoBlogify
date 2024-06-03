package post_image

import "time"

type PostImage struct {
	Id           int
	Post_Id      int
	Image_1      string
	Image_Name_1 string
	Image_2      string
	Image_Name_2 string
	Image_3      string
	Image_Name_3 string
	Created_At   time.Time
	Updated_At   time.Time
}
