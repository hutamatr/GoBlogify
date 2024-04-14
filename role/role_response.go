package role

import (
	"time"
)

type RoleResponse struct {
	Id         int       `json:"id"`
	Name       string    `json:"name"`
	Created_At time.Time `json:"created_at"`
	Updated_At time.Time `json:"updated_at"`
}

func ToRoleResponse(role Role) RoleResponse {
	return RoleResponse{
		Id:         role.Id,
		Name:       role.Name,
		Created_At: role.Created_At,
		Updated_At: role.Updated_At,
	}
}
