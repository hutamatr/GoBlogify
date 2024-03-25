package test

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	controllersCategory "github.com/hutamatr/GoBlogify/controllers/category"
	controllersComment "github.com/hutamatr/GoBlogify/controllers/comment"
	controllersPost "github.com/hutamatr/GoBlogify/controllers/post"
	controllersRole "github.com/hutamatr/GoBlogify/controllers/role"
	controllersUser "github.com/hutamatr/GoBlogify/controllers/user"
	"github.com/hutamatr/GoBlogify/helpers"
	repositoriesCategory "github.com/hutamatr/GoBlogify/repositories/category"
	repositoriesComment "github.com/hutamatr/GoBlogify/repositories/comment"
	repositoriesPost "github.com/hutamatr/GoBlogify/repositories/post"
	repositoriesRole "github.com/hutamatr/GoBlogify/repositories/role"
	repositoriesUser "github.com/hutamatr/GoBlogify/repositories/user"
	"github.com/hutamatr/GoBlogify/routes"
	servicesCategory "github.com/hutamatr/GoBlogify/services/category"
	servicesComment "github.com/hutamatr/GoBlogify/services/comment"
	servicesPost "github.com/hutamatr/GoBlogify/services/post"
	servicesRole "github.com/hutamatr/GoBlogify/services/role"
	servicesUser "github.com/hutamatr/GoBlogify/services/user"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load("../.env.test")
	helpers.PanicError(err)
}

func ConnectDBTest() *sql.DB {
	env := helpers.NewEnv()
	DBName := env.DB.DbName
	DBUsername := env.DB.Username
	DBPassword := env.DB.Password
	DBPort := env.DB.Port
	Host := env.DB.Host

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", DBUsername, DBPassword, Host, DBPort, DBName))
	helpers.PanicError(err)

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(1 * time.Hour)

	return db
}

func DeleteDBTest(db *sql.DB) {
	_, err := db.Exec("DELETE FROM comment")
	helpers.PanicError(err)
	_, err = db.Exec("DELETE FROM post")
	helpers.PanicError(err)
	_, err = db.Exec("DELETE FROM category")
	helpers.PanicError(err)
	_, err = db.Exec("DELETE FROM user")
	helpers.PanicError(err)
	_, err = db.Exec("DELETE FROM role")
	helpers.PanicError(err)
}

func SetupRouterTest(db *sql.DB) http.Handler {
	validator := validator.New()

	roleRepository := repositoriesRole.NewRoleRepository()
	roleService := servicesRole.NewRoleService(roleRepository, db, validator)
	roleController := controllersRole.NewRoleController(roleService)

	PostRepository := repositoriesPost.NewPostRepository()
	PostService := servicesPost.NewPostService(PostRepository, db, validator)
	PostController := controllersPost.NewPostController(PostService)

	categoryRepository := repositoriesCategory.NewCategoryRepository()
	categoryService := servicesCategory.NewCategoryService(categoryRepository, db, validator)
	categoryController := controllersCategory.NewCategoryController(categoryService)

	commentRepository := repositoriesComment.NewCommentRepository()
	commentService := servicesComment.NewCommentService(commentRepository, db, validator)
	commentController := controllersComment.NewCommentController(commentService)

	userRepository := repositoriesUser.NewUserRepository()
	userService := servicesUser.NewUserService(userRepository, roleRepository, db, validator)
	UserController := controllersUser.NewUserController(userService)

	router := routes.Router(&routes.RouterControllers{
		PostController:     PostController,
		CategoryController: categoryController,
		RoleController:     roleController,
		UserController:     UserController,
		CommentController:  commentController,
	})

	return router
}
