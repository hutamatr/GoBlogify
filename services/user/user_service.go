package services

import (
	"context"

	"github.com/hutamatr/go-blog-api/model/web"
)

type UserService interface {
	SignUp(ctx context.Context, request web.UserCreateRequest) (web.UserResponse, string, string)
	SignIn(ctx context.Context, request web.UserLoginRequest) (web.UserResponse, string, string)
	FindAll(ctx context.Context) []web.UserResponse
	FindById(ctx context.Context, userId int) web.UserResponse
	Update(ctx context.Context, request web.UserUpdateRequest) web.UserResponse
	Delete(ctx context.Context, userId int)
}
