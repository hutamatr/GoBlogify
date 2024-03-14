package routes

import (
	controllersA "github.com/hutamatr/go-blog-api/controllers/article"
	controllersC "github.com/hutamatr/go-blog-api/controllers/category"
	"github.com/hutamatr/go-blog-api/exception"
	"github.com/julienschmidt/httprouter"
)

type RouterControllers struct {
	ArticleController  controllersA.ArticleController
	CategoryController controllersC.CategoryController
}

func Router(routerControllers *RouterControllers) *httprouter.Router {
	router := httprouter.New()

	router.GET("/api/article", routerControllers.ArticleController.FindAllArticle)
	router.GET("/api/article/:articleId", routerControllers.ArticleController.FindByIdArticle)
	router.POST("/api/article", routerControllers.ArticleController.CreateArticle)
	router.PUT("/api/article/:articleId", routerControllers.ArticleController.UpdateArticle)
	router.DELETE("/api/article/:articleId", routerControllers.ArticleController.DeleteArticle)

	router.GET("/api/category", routerControllers.CategoryController.FindAllCategory)
	router.GET("/api/category/:categoryId", routerControllers.CategoryController.FindByIdCategory)
	router.POST("/api/category", routerControllers.CategoryController.CreateCategory)
	router.PUT("/api/category/:categoryId", routerControllers.CategoryController.UpdateCategory)
	router.DELETE("/api/category/:categoryId", routerControllers.CategoryController.DeleteCategory)

	router.PanicHandler = exception.ErrorHandler

	return router
}
