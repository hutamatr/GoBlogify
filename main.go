package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/hutamatr/go-blog-api/app"
	controllersA "github.com/hutamatr/go-blog-api/controllers/article"
	controllersC "github.com/hutamatr/go-blog-api/controllers/category"
	controllersU "github.com/hutamatr/go-blog-api/controllers/user"
	"github.com/hutamatr/go-blog-api/helpers"
	repositoriesA "github.com/hutamatr/go-blog-api/repositories/article"
	repositoriesC "github.com/hutamatr/go-blog-api/repositories/category"
	repositoriesU "github.com/hutamatr/go-blog-api/repositories/user"
	"github.com/hutamatr/go-blog-api/routes"
	servicesA "github.com/hutamatr/go-blog-api/services/article"
	servicesC "github.com/hutamatr/go-blog-api/services/category"
	servicesU "github.com/hutamatr/go-blog-api/services/user"

	"github.com/joho/godotenv"
)

func init() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	godotenv.Load(".env." + env)
	if env != "test" {
		godotenv.Load(".env")
	}
	godotenv.Load(".env." + env)
	godotenv.Load()
}

func main() {
	db := app.ConnectDB()
	validator := validator.New(validator.WithRequiredStructEnabled())

	userRepository := repositoriesU.NewUserRepository()
	userService := servicesU.NewUserService(userRepository, db, validator)
	UserController := controllersU.NewUserController(userService)

	articleRepository := repositoriesA.NewArticleRepository()
	articleService := servicesA.NewArticleService(articleRepository, db, validator)
	articleController := controllersA.NewArticleController(articleService)

	categoryRepository := repositoriesC.NewCategoryRepository()
	categoryService := servicesC.NewCategoryService(categoryRepository, db, validator)
	categoryController := controllersC.NewCategoryController(categoryService)

	router := routes.Router(&routes.RouterControllers{
		UserController:     UserController,
		ArticleController:  articleController,
		CategoryController: categoryController,
	})

	cors := helpers.Cors()

	handler := cors.Handler(router)

	server := http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	helpers.ServerRunningText()

	log.Fatal(server.ListenAndServe())
}
