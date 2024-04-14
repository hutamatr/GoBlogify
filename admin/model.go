package admin

import "time"

type Admin struct {
	Id         int
	Role_Id    int
	Username   string
	Email      string
	Password   string
	First_Name string
	Last_Name  string
	Created_At time.Time
	Updated_At time.Time
	Deleted_At time.Time
}
