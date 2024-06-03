package post

import (
	"time"

	"github.com/hutamatr/GoBlogify/category"
	"github.com/hutamatr/GoBlogify/post_image"
	"github.com/hutamatr/GoBlogify/user"
)

type Post struct {
	Id          int
	User_Id     int
	Category_Id int
	Title       string
	Post_Body   string
	Published   bool
	Deleted     bool
	Created_At  time.Time
	Updated_At  time.Time
	Deleted_At  time.Time
}

type PostCreateOrUpdate struct {
	Id          int
	User_id     int
	Category_Id int
	Title       string
	Post_Body   string
	Published   bool
	Deleted     bool
}

type PostJoin struct {
	Id         int
	Title      string
	Post_Body  string
	Published  bool
	Deleted    bool
	Created_At time.Time
	Updated_At time.Time
	Deleted_At time.Time
	User       user.UserJoin
	Category   category.Category
	Images     post_image.PostImage
}

type PostJoinFollowed struct {
	Id          int
	User_Id     int
	Category_Id int
	Title       string
	Post_Body   string
	Published   bool
	Deleted     bool
	Created_At  time.Time
	Updated_At  time.Time
	Deleted_At  time.Time
	User        user.UserJoin
	Images      post_image.PostImage
}
