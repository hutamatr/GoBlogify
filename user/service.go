package user

import (
	"context"
	"database/sql"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/hutamatr/GoBlogify/exception"
	"github.com/hutamatr/GoBlogify/helpers"
	"github.com/hutamatr/GoBlogify/role"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	SignUp(ctx context.Context, request UserCreateRequest) (UserResponse, string, string)
	SignIn(ctx context.Context, request UserLoginRequest) (UserResponse, string, string)
	FindAll(ctx context.Context, isAdmin bool) []UserResponse
	FindById(ctx context.Context, userId int) UserResponse
	Update(ctx context.Context, request UserUpdateRequest) UserResponse
	Delete(ctx context.Context, userId int)
}

type UserServiceImpl struct {
	userRepository UserRepository
	roleRepository role.RoleRepository
	DB             *sql.DB
	Validator      *validator.Validate
}

func NewUserService(userRepository UserRepository, roleRepository role.RoleRepository, db *sql.DB, validator *validator.Validate) UserService {
	return &UserServiceImpl{
		userRepository: userRepository,
		roleRepository: roleRepository,
		DB:             db,
		Validator:      validator,
	}
}

func (service *UserServiceImpl) SignUp(ctx context.Context, request UserCreateRequest) (UserResponse, string, string) {
	env := helpers.NewEnv()
	appEnv := env.App.AppEnv
	accessTokenSecret := env.SecretToken.AccessSecret
	refreshTokenSecret := env.SecretToken.RefreshSecret

	err := service.Validator.Struct(request)
	helpers.PanicError(err, "invalid request")

	tx, err := service.DB.Begin()
	helpers.PanicError(err, "failed to begin transaction")
	defer helpers.TxRollbackCommit(tx)

	user := service.userRepository.FindOne(ctx, tx, 0, request.Email)

	if user.Email == request.Email {
		panic(exception.NewBadRequestError("email already exist"))
	}

	userRole := service.roleRepository.FindByName(ctx, tx, "user")

	if userRole.Name != "user" {
		newRole := role.Role{
			Name: "user",
		}
		userRole = service.roleRepository.Save(ctx, tx, newRole)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)

	helpers.PanicError(err, "failed to hash password")

	newUser := User{
		Username: request.Username,
		Email:    request.Email,
		Password: string(hashedPassword),
		Role_Id:  userRole.Id,
	}

	createdUser := service.userRepository.Save(ctx, tx, newUser)

	accessTokenExpired := helpers.AccessTokenDuration(appEnv)

	accessToken, err := helpers.GenerateToken(createdUser.Id, accessTokenExpired, accessTokenSecret)
	helpers.PanicError(err, "failed to generate access token")

	refreshToken, err := helpers.GenerateToken(createdUser.Id, 168*time.Hour, refreshTokenSecret)
	helpers.PanicError(err, "failed to generate refresh token")

	return ToUserResponse(createdUser), accessToken, refreshToken
}

func (service *UserServiceImpl) SignIn(ctx context.Context, request UserLoginRequest) (UserResponse, string, string) {
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

	user := service.userRepository.FindOne(ctx, tx, 0, request.Email)

	accessTokenExpired := helpers.AccessTokenDuration(appEnv)

	accessToken, err := helpers.GenerateToken(user.Id, accessTokenExpired, accessTokenSecret)
	helpers.PanicError(err, "failed to generate access token")

	refreshToken, err := helpers.GenerateToken(user.Id, 168*time.Hour, refreshTokenSecret)
	helpers.PanicError(err, "failed to generate refresh token")

	return ToUserResponse(user), accessToken, refreshToken
}

func (service *UserServiceImpl) FindById(ctx context.Context, userId int) UserResponse {
	tx, err := service.DB.Begin()
	helpers.PanicError(err, "failed to begin transaction")
	defer helpers.TxRollbackCommit(tx)

	user := service.userRepository.FindOne(ctx, tx, userId, "")

	return ToUserResponse(user)
}

func (service *UserServiceImpl) FindAll(ctx context.Context, isAdmin bool) []UserResponse {
	if !isAdmin {
		panic(exception.NewBadRequestError("only admin can get all users"))
	}

	tx, err := service.DB.Begin()
	helpers.PanicError(err, "failed to begin transaction")
	defer helpers.TxRollbackCommit(tx)

	users := service.userRepository.FindAll(ctx, tx)

	var usersData []UserResponse

	for _, user := range users {
		usersData = append(usersData, ToUserResponse(user))
	}

	return usersData
}

func (service *UserServiceImpl) Update(ctx context.Context, request UserUpdateRequest) UserResponse {
	err := service.Validator.Struct(request)
	helpers.PanicError(err, "invalid request")

	tx, err := service.DB.Begin()
	helpers.PanicError(err, "failed to begin transaction")
	defer helpers.TxRollbackCommit(tx)

	user := service.userRepository.FindOne(ctx, tx, request.Id, "")

	user.Username = request.Username
	user.First_Name = request.First_Name
	user.Last_Name = request.Last_Name

	updatedUser := service.userRepository.Update(ctx, tx, user)

	return ToUserResponse(updatedUser)
}

func (service *UserServiceImpl) Delete(ctx context.Context, userId int) {
	tx, err := service.DB.Begin()
	helpers.PanicError(err, "failed to begin transaction")
	defer helpers.TxRollbackCommit(tx)

	user := service.userRepository.FindOne(ctx, tx, userId, "")

	if user.Id <= 0 {
		panic(exception.NewNotFoundError("user not found"))
	}

	service.userRepository.Delete(ctx, tx, user.Id)
}
