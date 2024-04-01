package services

import (
	"context"
	"database/sql"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/hutamatr/GoBlogify/exception"
	"github.com/hutamatr/GoBlogify/helpers"
	"github.com/hutamatr/GoBlogify/model/domain"
	"github.com/hutamatr/GoBlogify/model/web"
	repositoriesRole "github.com/hutamatr/GoBlogify/repositories/role"
	repositoriesUser "github.com/hutamatr/GoBlogify/repositories/user"
	"golang.org/x/crypto/bcrypt"
)

type AdminServiceImpl struct {
	userRepository repositoriesUser.UserRepository
	roleRepository repositoriesRole.RoleRepository
	DB             *sql.DB
	Validator      *validator.Validate
}

func NewAdminService(userRepository repositoriesUser.UserRepository, roleRepository repositoriesRole.RoleRepository, DB *sql.DB, Validator *validator.Validate) AdminService {
	return &AdminServiceImpl{
		userRepository: userRepository,
		roleRepository: roleRepository,
		DB:             DB,
		Validator:      Validator,
	}
}

func (service *AdminServiceImpl) SignUpAdmin(ctx context.Context, request web.AdminCreateRequest) (web.AdminResponse, string, string) {
	env := helpers.NewEnv()
	appEnv := env.App.AppEnv
	adminCode := env.Auth.AdminCode
	accessTokenSecret := env.SecretToken.AccessSecret
	refreshTokenSecret := env.SecretToken.RefreshSecret

	err := service.Validator.Struct(request)
	helpers.PanicError(err)

	tx, err := service.DB.Begin()
	helpers.PanicError(err)
	defer helpers.TxRollbackCommit(tx)

	admin := service.userRepository.FindOne(ctx, tx, 0, request.Email)

	if admin.Email == request.Email {
		panic(exception.NewBadRequestError("email already exist"))
	}

	if adminCode != request.Admin_Code {
		panic(exception.NewBadRequestError("invalid admin code"))
	}

	role := service.roleRepository.FindByName(ctx, tx, "admin")

	if role.Name != "admin" {
		newRole := domain.Role{
			Name: "admin",
		}
		role = service.roleRepository.Save(ctx, tx, newRole)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)

	helpers.PanicError(err)

	newAdmin := domain.User{
		Username: request.Username,
		Email:    request.Email,
		Password: string(hashedPassword),
		Role_Id:  role.Id,
	}

	createdAdmin := service.userRepository.Save(ctx, tx, newAdmin)

	accessTokenExpired := helpers.AccessTokenDuration(appEnv)

	accessToken, err := helpers.GenerateToken(createdAdmin.Id, accessTokenExpired, accessTokenSecret)
	helpers.PanicError(err)

	refreshToken, err := helpers.GenerateToken(createdAdmin.Id, 168*time.Hour, refreshTokenSecret)
	helpers.PanicError(err)

	return web.ToAdminResponse(createdAdmin), accessToken, refreshToken
}

func (service *AdminServiceImpl) SignInAdmin(ctx context.Context, request web.AdminLoginRequest) (web.AdminResponse, string, string) {
	env := helpers.NewEnv()
	appEnv := env.App.AppEnv
	accessTokenSecret := env.SecretToken.AccessSecret
	refreshTokenSecret := env.SecretToken.RefreshSecret

	err := service.Validator.Struct(request)
	helpers.PanicError(err)

	tx, err := service.DB.Begin()
	helpers.PanicError(err)
	defer helpers.TxRollbackCommit(tx)

	password := service.userRepository.FindPassword(ctx, tx, request.Email)

	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(request.Password))
	if err != nil {
		panic(exception.NewBadRequestError("invalid email or password"))
	}

	admin := service.userRepository.FindOne(ctx, tx, 0, request.Email)

	accessTokenExpired := helpers.AccessTokenDuration(appEnv)

	accessToken, err := helpers.GenerateToken(admin.Id, accessTokenExpired, accessTokenSecret)
	helpers.PanicError(err)

	refreshToken, err := helpers.GenerateToken(admin.Id, 168*time.Hour, refreshTokenSecret)
	helpers.PanicError(err)

	return web.ToAdminResponse(admin), accessToken, refreshToken
}
