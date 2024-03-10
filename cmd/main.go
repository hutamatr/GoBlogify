package main

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	blog "github.com/hutamatr/go-blog-api/cmd/go_blog"
	"github.com/hutamatr/go-blog-api/cmd/go_blog/helper"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load("../.env")
	helper.PanicError(err)
}

func main() {
	db := blog.ConnectDB()
	validator := validator.New()

	articleRepository := blog.NewArticleRepository()
	articleService := blog.NewArticleService(articleRepository, db, validator)
	articleController := blog.NewArticleController(articleService)

	categoryRepository := blog.NewCategoryRepository()
	categoryService := blog.NewCategoryService(categoryRepository, db, validator)
	categoryController := blog.NewCategoryController(categoryService)

	router := blog.Router(articleController, categoryController)

	cors := helper.Cors()

	handler := cors.Handler(router)

	server := http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	err := server.ListenAndServe()
	helper.PanicError(err)
}
