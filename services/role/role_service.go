package services

import (
	"context"

	"github.com/hutamatr/go-blog-api/model/web"
)

type RoleService interface {
	Create(ctx context.Context, request web.RoleCreateRequest) web.RoleResponse
	FindAll(ctx context.Context) []web.RoleResponse
	FindById(ctx context.Context, roleId int) web.RoleResponse
	Update(ctx context.Context, request web.RoleUpdateRequest) web.RoleResponse
	Delete(ctx context.Context, roleId int)
}
