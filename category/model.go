package category

import "time"

type Category struct {
	Id         int
	Name       string
	Created_At time.Time
	Updated_At time.Time
}

type CategoryCreateOrUpdate struct {
	Id   int
	Name string
}
