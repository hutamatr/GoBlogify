package test

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	controllersA "github.com/hutamatr/go-blog-api/controllers/article"
	controllersC "github.com/hutamatr/go-blog-api/controllers/category"
	controllersR "github.com/hutamatr/go-blog-api/controllers/role"
	controllersU "github.com/hutamatr/go-blog-api/controllers/user"
	"github.com/hutamatr/go-blog-api/helpers"
	repositoriesA "github.com/hutamatr/go-blog-api/repositories/article"
	repositoriesC "github.com/hutamatr/go-blog-api/repositories/category"
	repositoriesR "github.com/hutamatr/go-blog-api/repositories/role"
	repositoriesU "github.com/hutamatr/go-blog-api/repositories/user"
	"github.com/hutamatr/go-blog-api/routes"
	servicesA "github.com/hutamatr/go-blog-api/services/article"
	servicesC "github.com/hutamatr/go-blog-api/services/category"
	servicesR "github.com/hutamatr/go-blog-api/services/role"
	servicesU "github.com/hutamatr/go-blog-api/services/user"
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
	_, err := db.Exec("DELETE FROM article")
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

	roleRepository := repositoriesR.NewRoleRepository()
	roleService := servicesR.NewRoleService(roleRepository, db, validator)
	roleController := controllersR.NewRoleController(roleService)

	articleRepository := repositoriesA.NewArticleRepository()
	articleService := servicesA.NewArticleService(articleRepository, db, validator)
	articleController := controllersA.NewArticleController(articleService)

	categoryRepository := repositoriesC.NewCategoryRepository()
	categoryService := servicesC.NewCategoryService(categoryRepository, db, validator)
	categoryController := controllersC.NewCategoryController(categoryService)

	userRepository := repositoriesU.NewUserRepository()
	userService := servicesU.NewUserService(userRepository, roleRepository, db, validator)
	UserController := controllersU.NewUserController(userService)

	router := routes.Router(&routes.RouterControllers{
		ArticleController:  articleController,
		CategoryController: categoryController,
		RoleController:     roleController,
		UserController:     UserController,
	})

	return router
}
