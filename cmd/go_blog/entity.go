package go_blog

import "time"

// type Article struct {
// 	Id          int
// 	Title       string
// 	Body        string
// 	Author      string
// 	Category_Id int
// 	Published   bool
// 	Deleted     bool
// 	Created_At  time.Time
// 	Updated_At  time.Time
// 	Deleted_At  time.Time
// }

type Category struct {
	Id         int
	Name       string
	Created_At time.Time
	Updated_At time.Time
}
