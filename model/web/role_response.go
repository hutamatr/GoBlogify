package web

import "time"

type RoleResponse struct {
	Id         int       `json:"id"`
	Name       string    `json:"name"`
	Created_At time.Time `json:"created_at"`
	Updated_At time.Time `json:"updated_at"`
}
