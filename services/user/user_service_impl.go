package services

import (
	"context"
	"database/sql"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/hutamatr/go-blog-api/exception"
	"github.com/hutamatr/go-blog-api/helpers"
	"github.com/hutamatr/go-blog-api/model/domain"
	"github.com/hutamatr/go-blog-api/model/web"
	repositoriesU "github.com/hutamatr/go-blog-api/repositories/user"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl struct {
	userRepository repositoriesU.UserRepository
	DB             *sql.DB
	Validator      *validator.Validate
}

var accessTokenSecret = os.Getenv("ACCESS_TOKEN_SECRET")
var refreshTokenSecret = os.Getenv("REFRESH_TOKEN_SECRET")

func NewUserService(userRepository repositoriesU.UserRepository, db *sql.DB, validator *validator.Validate) UserService {
	return &UserServiceImpl{
		userRepository: userRepository,
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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)

	helpers.PanicError(err)

	newUser := domain.User{
		Username: request.Username,
		Email:    request.Email,
		Password: string(hashedPassword),
		Role_Id:  1, // temporary
	}

	newUser = service.userRepository.Save(ctx, tx, newUser)

	accessToken, err := helpers.GenerateToken(user.Username, 5*time.Minute, accessTokenSecret)
	helpers.PanicError(err)

	refreshToken, err := helpers.GenerateToken(user.Username, 168*time.Hour, refreshTokenSecret)
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

	accessToken, err := helpers.GenerateToken(user.Username, 5*time.Minute, accessTokenSecret)
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
	user.Email = request.Email
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
