package admin

import (
	"context"
	"database/sql"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/hutamatr/GoBlogify/exception"
	"github.com/hutamatr/GoBlogify/helpers"
	"github.com/hutamatr/GoBlogify/role"
	"github.com/hutamatr/GoBlogify/user"
	"golang.org/x/crypto/bcrypt"
)

type AdminService interface {
	SignUpAdmin(ctx context.Context, request AdminCreateRequest) (AdminResponse, string, string)
	SignInAdmin(ctx context.Context, request AdminLoginRequest) (AdminResponse, string, string)
}

type AdminServiceImpl struct {
	userRepository user.UserRepository
	roleRepository role.RoleRepository
	DB             *sql.DB
	Validator      *validator.Validate
}

func NewAdminService(userRepository user.UserRepository, roleRepository role.RoleRepository, DB *sql.DB, Validator *validator.Validate) AdminService {
	return &AdminServiceImpl{
		userRepository: userRepository,
		roleRepository: roleRepository,
		DB:             DB,
		Validator:      Validator,
	}
}

func (service *AdminServiceImpl) SignUpAdmin(ctx context.Context, request AdminCreateRequest) (AdminResponse, string, string) {
	env := helpers.NewEnv()
	appEnv := env.App.AppEnv
	adminCode := env.Auth.AdminCode
	accessTokenSecret := env.SecretToken.AccessSecret
	refreshTokenSecret := env.SecretToken.RefreshSecret

	err := service.Validator.Struct(request)
	helpers.PanicError(err, "invalid request")

	tx, err := service.DB.Begin()
	helpers.PanicError(err, "failed to begin transaction")
	defer helpers.TxRollbackCommit(tx)

	admin := service.userRepository.FindOne(ctx, tx, 0, request.Email)

	if admin.Email == request.Email {
		panic(exception.NewBadRequestError("email already exist"))
	}

	if adminCode != request.Admin_Code {
		panic(exception.NewBadRequestError("invalid admin code"))
	}

	adminRole := service.roleRepository.FindByName(ctx, tx, "admin")

	if adminRole.Name != "admin" {
		newRole := role.Role{
			Name: "admin",
		}
		adminRole = service.roleRepository.Save(ctx, tx, newRole)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)

	helpers.PanicError(err, "failed to hash password")

	newAdmin := user.User{
		Username: request.Username,
		Email:    request.Email,
		Password: string(hashedPassword),
		Role_Id:  adminRole.Id,
	}

	createdAdmin := service.userRepository.Save(ctx, tx, newAdmin)

	accessTokenExpired := helpers.AccessTokenDuration(appEnv)

	accessToken, err := helpers.GenerateToken(createdAdmin.Id, accessTokenExpired, accessTokenSecret)
	helpers.PanicError(err, "failed to generate access token")

	refreshToken, err := helpers.GenerateToken(createdAdmin.Id, 168*time.Hour, refreshTokenSecret)
	helpers.PanicError(err, "failed to generate refresh token")

	return ToAdminResponse(createdAdmin), accessToken, refreshToken
}

func (service *AdminServiceImpl) SignInAdmin(ctx context.Context, request AdminLoginRequest) (AdminResponse, string, string) {
	env := helpers.NewEnv()
	appEnv := env.App.AppEnv
	accessTokenSecret := env.SecretToken.AccessSecret
	refreshTokenSecret := env.SecretToken.RefreshSecret

	err := service.Validator.Struct(request)
	helpers.PanicError(err, "invalid request")

	tx, err := service.DB.Begin()
	helpers.PanicError(err, "failed to begin transaction")
	defer helpers.TxRollbackCommit(tx)

	password := service.userRepository.FindPassword(ctx, tx, request.Email)

	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(request.Password))
	if err != nil {
		panic(exception.NewBadRequestError("invalid email or password"))
	}

	admin := service.userRepository.FindOne(ctx, tx, 0, request.Email)

	accessTokenExpired := helpers.AccessTokenDuration(appEnv)

	accessToken, err := helpers.GenerateToken(admin.Id, accessTokenExpired, accessTokenSecret)
	helpers.PanicError(err, "failed to generate access token")

	refreshToken, err := helpers.GenerateToken(admin.Id, 168*time.Hour, refreshTokenSecret)
	helpers.PanicError(err, "failed to generate refresh token")

	return ToAdminResponse(admin), accessToken, refreshToken
}
