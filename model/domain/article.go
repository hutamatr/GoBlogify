package domain

import "time"

type Article struct {
	Id          int
	Category_Id int
	Title       string
	Body        string
	Author      string
	Published   bool
	Deleted     bool
	Created_At  time.Time
	Updated_At  time.Time
	Deleted_At  time.Time
}

type ArticleCreateOrUpdate struct {
	Id          int
	Category_Id int
	Title       string
	Body        string
	Author      string
	Published   bool
	Deleted     bool
}

type ArticleJoinCategory struct {
	Id         int
	Title      string
	Body       string
	Author     string
	Published  bool
	Deleted    bool
	Created_At time.Time
	Updated_At time.Time
	Deleted_At time.Time
	Category   Category
}
