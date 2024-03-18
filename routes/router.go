package routes

import (
	"net/http"

	controllersA "github.com/hutamatr/go-blog-api/controllers/article"
	controllersC "github.com/hutamatr/go-blog-api/controllers/category"
	controllersR "github.com/hutamatr/go-blog-api/controllers/role"
	controllersU "github.com/hutamatr/go-blog-api/controllers/user"
	"github.com/hutamatr/go-blog-api/exception"
	"github.com/hutamatr/go-blog-api/helpers"
	"github.com/hutamatr/go-blog-api/model/web"
	"github.com/julienschmidt/httprouter"
)

type RouterControllers struct {
	UserController     controllersU.UserController
	ArticleController  controllersA.ArticleController
	CategoryController controllersC.CategoryController
	RoleController     controllersR.RoleController
}

func Router(routerControllers *RouterControllers) *httprouter.Router {
	router := httprouter.New()

	router.POST("/api/signup", routerControllers.UserController.CreateUser)
	router.POST("/api/signin", routerControllers.UserController.SignInUser)
	router.GET("/api/user", routerControllers.UserController.FindAllUser)
	router.GET("/api/user/:userId", routerControllers.UserController.FindByIdUser)
	router.PUT("/api/user/:userId", routerControllers.UserController.UpdateUser)
	router.DELETE("/api/user/:userId", routerControllers.UserController.DeleteUser)
	router.GET("/api/refresh", routerControllers.UserController.GetRefreshToken)

	router.POST("/api/role", routerControllers.RoleController.CreateRole)
	router.GET("/api/role", routerControllers.RoleController.FindAllRole)
	router.GET("/api/role/:roleId", routerControllers.RoleController.FindRoleById)
	router.PUT("/api/role/:roleId", routerControllers.RoleController.UpdateRole)
	router.DELETE("/api/role/:roleId", routerControllers.RoleController.DeleteRole)

	router.POST("/api/article", routerControllers.ArticleController.CreateArticle)
	router.GET("/api/article", routerControllers.ArticleController.FindAllArticle)
	router.GET("/api/article/:articleId", routerControllers.ArticleController.FindByIdArticle)
	router.PUT("/api/article/:articleId", routerControllers.ArticleController.UpdateArticle)
	router.DELETE("/api/article/:articleId", routerControllers.ArticleController.DeleteArticle)

	router.POST("/api/category", routerControllers.CategoryController.CreateCategory)
	router.GET("/api/category", routerControllers.CategoryController.FindAllCategory)
	router.GET("/api/category/:categoryId", routerControllers.CategoryController.FindByIdCategory)
	router.PUT("/api/category/:categoryId", routerControllers.CategoryController.UpdateCategory)
	router.DELETE("/api/category/:categoryId", routerControllers.CategoryController.DeleteCategory)

	router.NotFound = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusNotFound)
		userResponse := web.ResponseJSON{
			Code:   http.StatusNotFound,
			Status: "NOT FOUND",
		}

		helpers.EncodeJSONFromResponse(writer, userResponse)
	})

	router.PanicHandler = exception.ErrorHandler

	return router
}
