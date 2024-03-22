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

type UserServiceImpl struct {
	userRepository repositoriesUser.UserRepository
	roleRepository repositoriesRole.RoleRepository
	DB             *sql.DB
	Validator      *validator.Validate
}

var env = helpers.NewEnv()
var appEnv = env.App.AppEnv
var accessTokenSecret = env.SecretToken.AccessSecret
var refreshTokenSecret = env.SecretToken.RefreshSecret

func NewUserService(userRepository repositoriesUser.UserRepository, roleRepository repositoriesRole.RoleRepository, db *sql.DB, validator *validator.Validate) UserService {
	return &UserServiceImpl{
		userRepository: userRepository,
		roleRepository: roleRepository,
		DB:             db,
		Validator:      validator,
	}
}

func (service *UserServiceImpl) SignUp(ctx context.Context, request web.UserCreateRequest) (web.UserResponse, string, string) {
	err := service.Validator.Struct(request)
	helpers.PanicError(err)

	tx, err := service.DB.Begin()
	helpers.PanicError(err)
	defer helpers.TxRollbackCommit(tx)

	user := service.userRepository.FindOne(ctx, tx, 0, request.Email)

	if user.Email == request.Email {
		panic(exception.NewBadRequestError("email already exist"))
	}

	role := service.roleRepository.FindByName(ctx, tx, "user")

	if role.Name != "user" {
		newRole := domain.Role{
			Name: "user",
		}
		role = service.roleRepository.Save(ctx, tx, newRole)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)

	helpers.PanicError(err)

	newUser := domain.User{
		Username: request.Username,
		Email:    request.Email,
		Password: string(hashedPassword),
		Role_Id:  role.Id,
	}

	newUser = service.userRepository.Save(ctx, tx, newUser)

	accessTokenExpired := helpers.AccessTokenDuration(appEnv)

	accessToken, err := helpers.GenerateToken(request.Username, accessTokenExpired, accessTokenSecret)
	helpers.PanicError(err)

	refreshToken, err := helpers.GenerateToken(request.Username, 168*time.Hour, refreshTokenSecret)
	helpers.PanicError(err)

	return web.ToUserResponse(newUser), accessToken, refreshToken
}

func (service *UserServiceImpl) SignIn(ctx context.Context, request web.UserLoginRequest) (web.UserResponse, string, string) {
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

	user := service.userRepository.FindOne(ctx, tx, 0, request.Email)

	accessTokenExpired := helpers.AccessTokenDuration(appEnv)

	accessToken, err := helpers.GenerateToken(user.Username, accessTokenExpired, accessTokenSecret)
	helpers.PanicError(err)

	refreshToken, err := helpers.GenerateToken(user.Username, 168*time.Hour, refreshTokenSecret)
	helpers.PanicError(err)

	return web.ToUserResponse(user), accessToken, refreshToken
}

func (service *UserServiceImpl) FindById(ctx context.Context, userId int) web.UserResponse {
	tx, err := service.DB.Begin()
	helpers.PanicError(err)
	defer helpers.TxRollbackCommit(tx)

	user := service.userRepository.FindOne(ctx, tx, userId, "")

	return web.ToUserResponse(user)
}

func (service *UserServiceImpl) FindAll(ctx context.Context) []web.UserResponse {
	tx, err := service.DB.Begin()
	helpers.PanicError(err)
	defer helpers.TxRollbackCommit(tx)

	users := service.userRepository.FindAll(ctx, tx)

	var usersData []web.UserResponse

	for _, user := range users {
		usersData = append(usersData, web.ToUserResponse(user))
	}

	return usersData
}

func (service *UserServiceImpl) Update(ctx context.Context, request web.UserUpdateRequest) web.UserResponse {
	err := service.Validator.Struct(request)
	helpers.PanicError(err)

	tx, err := service.DB.Begin()
	helpers.PanicError(err)
	defer helpers.TxRollbackCommit(tx)

	user := service.userRepository.FindOne(ctx, tx, request.Id, "")

	user.Username = request.Username
	user.First_Name = request.First_Name
	user.Last_Name = request.Last_Name

	updatedUser := service.userRepository.Update(ctx, tx, user)

	return web.ToUserResponse(updatedUser)
}

func (service *UserServiceImpl) Delete(ctx context.Context, userId int) {
	tx, err := service.DB.Begin()
	helpers.PanicError(err)
	defer helpers.TxRollbackCommit(tx)

	user := service.userRepository.FindOne(ctx, tx, userId, "")

	if user.Id <= 0 {
		panic(exception.NewNotFoundError("user not found"))
	}

	service.userRepository.Delete(ctx, tx, user.Id)
}
