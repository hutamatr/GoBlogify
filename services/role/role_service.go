package services

import (
	"context"

	"github.com/hutamatr/GoBlogify/model/web"
)

type RoleService interface {
	Create(ctx context.Context, request web.RoleCreateRequest, isAdmin bool) web.RoleResponse
	FindAll(ctx context.Context, isAdmin bool) []web.RoleResponse
	FindById(ctx context.Context, roleId int, isAdmin bool) web.RoleResponse
	Update(ctx context.Context, request web.RoleUpdateRequest, isAdmin bool) web.RoleResponse
	Delete(ctx context.Context, roleId int, isAdmin bool)
}
