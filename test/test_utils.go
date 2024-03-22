package test

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	controllersC "github.com/hutamatr/GoBlogify/controllers/category"
	controllersP "github.com/hutamatr/GoBlogify/controllers/post"
	controllersR "github.com/hutamatr/GoBlogify/controllers/role"
	controllersU "github.com/hutamatr/GoBlogify/controllers/user"
	"github.com/hutamatr/GoBlogify/helpers"
	repositoriesC "github.com/hutamatr/GoBlogify/repositories/category"
	repositoriesP "github.com/hutamatr/GoBlogify/repositories/post"
	repositoriesR "github.com/hutamatr/GoBlogify/repositories/role"
	repositoriesU "github.com/hutamatr/GoBlogify/repositories/user"
	"github.com/hutamatr/GoBlogify/routes"
	servicesC "github.com/hutamatr/GoBlogify/services/category"
	servicesP "github.com/hutamatr/GoBlogify/services/post"
	servicesR "github.com/hutamatr/GoBlogify/services/role"
	servicesU "github.com/hutamatr/GoBlogify/services/user"
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
	_, err := db.Exec("DELETE FROM post")
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

	PostRepository := repositoriesP.NewPostRepository()
	PostService := servicesP.NewPostService(PostRepository, db, validator)
	PostController := controllersP.NewPostController(PostService)

	categoryRepository := repositoriesC.NewCategoryRepository()
	categoryService := servicesC.NewCategoryService(categoryRepository, db, validator)
	categoryController := controllersC.NewCategoryController(categoryService)

	userRepository := repositoriesU.NewUserRepository()
	userService := servicesU.NewUserService(userRepository, roleRepository, db, validator)
	UserController := controllersU.NewUserController(userService)

	router := routes.Router(&routes.RouterControllers{
		PostController:     PostController,
		CategoryController: categoryController,
		RoleController:     roleController,
		UserController:     UserController,
	})

	return router
}
