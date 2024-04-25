package post

import (
	"time"

	"github.com/hutamatr/GoBlogify/category"
	"github.com/hutamatr/GoBlogify/user"
)

type Post struct {
	Id          int
	User_Id     int
	Category_Id int
	Title       string
	Body        string
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
	Body        string
	Published   bool
	Deleted     bool
}

type PostJoin struct {
	Id         int
	Title      string
	Body       string
	Published  bool
	Deleted    bool
	Created_At time.Time
	Updated_At time.Time
	Deleted_At time.Time
	User       user.UserJoin
	Category   category.Category
}

type PostJoinFollowed struct {
	Id          int
	User_Id     int
	Category_Id int
	Title       string
	Body        string
	Published   bool
	Deleted     bool
	Created_At  time.Time
	Updated_At  time.Time
	Deleted_At  time.Time
	User        user.UserJoin
}
