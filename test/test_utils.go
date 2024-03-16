package test

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	controllersA "github.com/hutamatr/go-blog-api/controllers/article"
	controllersC "github.com/hutamatr/go-blog-api/controllers/category"
	"github.com/hutamatr/go-blog-api/helpers"
	repositoriesA "github.com/hutamatr/go-blog-api/repositories/article"
	repositoriesC "github.com/hutamatr/go-blog-api/repositories/category"
	"github.com/hutamatr/go-blog-api/routes"
	servicesA "github.com/hutamatr/go-blog-api/services/article"
	servicesC "github.com/hutamatr/go-blog-api/services/category"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load("../.env.test")
	helpers.PanicError(err)
}

func ConnectDBTest() *sql.DB {
	var DBName = os.Getenv("DB_NAME")
	var DBUsername = os.Getenv("DB_USERNAME")
	var DBPassword = os.Getenv("DB_PASSWORD")
	var DBPort = os.Getenv("DB_PORT")
	var Host = os.Getenv("HOST")

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
}

func SetupRouterTest(db *sql.DB) http.Handler {
	validator := validator.New()

	articleRepository := repositoriesA.NewArticleRepository()
	articleService := servicesA.NewArticleService(articleRepository, db, validator)
	articleController := controllersA.NewArticleController(articleService)

	categoryRepository := repositoriesC.NewCategoryRepository()
	categoryService := servicesC.NewCategoryService(categoryRepository, db, validator)
	categoryController := controllersC.NewCategoryController(categoryService)

	router := routes.Router(&routes.RouterControllers{
		ArticleController:  articleController,
		CategoryController: categoryController,
	})

	return router
}
