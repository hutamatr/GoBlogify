package comment

import (
	"time"

	"github.com/hutamatr/GoBlogify/user"
)

type Comment struct {
	Id         int
	Post_Id    int
	User_Id    int
	Content    string
	Created_At time.Time
	Updated_At time.Time
}

type CommentCreateOrUpdate struct {
	Id      int
	Post_Id int
	User_Id int
	Content string
}

type CommentJoin struct {
	Id         int
	Post_Id    int
	User_Id    int
	Content    string
	Created_At time.Time
	Updated_At time.Time
	User       user.User
}
