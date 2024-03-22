package domain

import "time"

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
	Content    string
	Post_Id    int
	User_Id    int
	Created_At time.Time
	Updated_At time.Time
	User       User
}
