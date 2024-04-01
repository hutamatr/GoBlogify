package web

import (
	"time"

	"github.com/hutamatr/GoBlogify/model/domain"
)

type AdminResponse struct {
	Id         int       `json:"id"`
	Role_Id    int       `json:"role_id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	First_Name string    `json:"first_name"`
	Last_Name  string    `json:"last_name"`
	Created_At time.Time `json:"created_at"`
	Updated_At time.Time `json:"updated_at"`
	Deleted_At time.Time `json:"deleted_at"`
}

func ToAdminResponse(user domain.User) AdminResponse {
	return AdminResponse{
		Id:         user.Id,
		Role_Id:    user.Role_Id,
		Username:   user.Username,
		Email:      user.Email,
		First_Name: user.First_Name,
		Last_Name:  user.Last_Name,
		Created_At: user.Created_At,
		Updated_At: user.Updated_At,
		Deleted_At: user.Deleted_At,
	}
}
