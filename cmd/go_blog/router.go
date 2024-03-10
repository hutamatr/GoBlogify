package go_blog

import "github.com/julienschmidt/httprouter"

func Router(articleController ArticleController, categoryController CategoryController) *httprouter.Router {

	router := httprouter.New()

	router.GET("/api/article", articleController.FindAllArticle)
	router.GET("/api/article/:articleId", articleController.FindByIdArticle)
	router.POST("/api/article", articleController.CreateArticle)
	router.PUT("/api/article/:articleId", articleController.UpdateArticle)
	router.DELETE("/api/article/:articleId", articleController.DeleteArticle)

	router.GET("/api/category", categoryController.FindAllCategory)
	router.GET("/api/category/:categoryId", categoryController.FindByIdCategory)
	router.POST("/api/category", categoryController.CreateCategory)
	router.PUT("/api/category/:categoryId", categoryController.UpdateCategory)
	router.DELETE("/api/category/:categoryId", categoryController.DeleteCategory)

	return router
}
