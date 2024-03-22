package web

import (
	"time"

	"github.com/hutamatr/GoBlogify/model/domain"
)

type CategoryResponse struct {
	Id         int       `json:"id"`
	Name       string    `json:"name"`
	Created_At time.Time `json:"created_at"`
	Updated_At time.Time `json:"updated_at"`
}

func ToCategoryResponse(category domain.Category) CategoryResponse {
	return CategoryResponse{
		Id:         category.Id,
		Name:       category.Name,
		Created_At: category.Created_At,
		Updated_At: category.Updated_At,
	}
}
