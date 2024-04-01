package services

import (
	"context"

	"github.com/hutamatr/GoBlogify/model/web"
)

type AdminService interface {
	SignUpAdmin(ctx context.Context, request web.AdminCreateRequest) (web.AdminResponse, string, string)
	SignInAdmin(ctx context.Context, request web.AdminLoginRequest) (web.AdminResponse, string, string)
}
