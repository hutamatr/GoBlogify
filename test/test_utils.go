package test

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hutamatr/GoBlogify/admin"
	"github.com/hutamatr/GoBlogify/category"
	"github.com/hutamatr/GoBlogify/comment"
	"github.com/hutamatr/GoBlogify/follow"
	"github.com/hutamatr/GoBlogify/post"
	"github.com/hutamatr/GoBlogify/role"
	"github.com/hutamatr/GoBlogify/routes"
	"github.com/hutamatr/GoBlogify/user"

	"github.com/hutamatr/GoBlogify/helpers"
	"github.com/hutamatr/GoBlogify/middleware"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load("../.env.test")
	helpers.PanicError(err, "failed to load .env.test")
}

func ConnectDBTest() *sql.DB {
	env := helpers.NewEnv()
	DBName := env.DB.DbName
	DBUsername := env.DB.Username
	DBPassword := env.DB.Password
	DBPort := env.DB.Port
	Host := env.DB.Host

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", DBUsername, DBPassword, Host, DBPort, DBName))
	helpers.PanicError(err, "failed to connect test database")

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(1 * time.Hour)

	return db
}

func DeleteDBTest(db *sql.DB) {
	_, err := db.Exec("DELETE FROM comment")
	helpers.PanicError(err, "failed to delete comment")
	_, err = db.Exec("DELETE FROM post")
	helpers.PanicError(err, "failed to delete post")
	_, err = db.Exec("DELETE FROM category")
	helpers.PanicError(err, "failed to delete category")
	_, err = db.Exec("DELETE FROM user")
	helpers.PanicError(err, "failed to delete user")
	_, err = db.Exec("DELETE FROM role")
	helpers.PanicError(err, "failed to delete role")
	_, err = db.Exec("DELETE FROM follow")
	helpers.PanicError(err, "failed to delete follow")
}

func SetupRouterTest(db *sql.DB) http.Handler {
	validator := validator.New()

	roleRepository := role.NewRoleRepository()
	roleService := role.NewRoleService(roleRepository, db, validator)
	roleController := role.NewRoleController(roleService)

	postRepository := post.NewPostRepository()
	postService := post.NewPostService(postRepository, db, validator)
	postController := post.NewPostController(postService)

	categoryRepository := category.NewCategoryRepository()
	categoryService := category.NewCategoryService(categoryRepository, db, validator)
	categoryController := category.NewCategoryController(categoryService)

	commentRepository := comment.NewCommentRepository()
	commentService := comment.NewCommentService(commentRepository, db, validator)
	commentController := comment.NewCommentController(commentService)

	userRepository := user.NewUserRepository()
	userService := user.NewUserService(userRepository, roleRepository, db, validator)
	userController := user.NewUserController(userService)

	adminService := admin.NewAdminService(userRepository, roleRepository, db, validator)
	adminController := admin.NewAdminControllerImpl(adminService)

	followRepository := follow.NewFollowRepository()
	followService := follow.NewFollowService(followRepository, db)
	followController := follow.NewFollowController(followService)

	router := routes.Router(&routes.RouterControllers{
		Admin:    adminController,
		User:     userController,
		Post:     postController,
		Category: categoryController,
		Role:     roleController,
		Comment:  commentController,
		Follow:   followController,
	})

	return middleware.NewAuthMiddleware(router)
}
