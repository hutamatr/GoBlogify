package role

import "time"

type Role struct {
	Id         int
	Name       string
	Created_At time.Time
	Updated_At time.Time
}
