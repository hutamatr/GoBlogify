package routes

import (
	"net/http"

	controllersCategory "github.com/hutamatr/GoBlogify/controllers/category"
	controllersComment "github.com/hutamatr/GoBlogify/controllers/comment"
	controllersPost "github.com/hutamatr/GoBlogify/controllers/post"
	controllersRole "github.com/hutamatr/GoBlogify/controllers/role"
	controllersUser "github.com/hutamatr/GoBlogify/controllers/user"
	"github.com/hutamatr/GoBlogify/exception"
	"github.com/hutamatr/GoBlogify/helpers"
	"github.com/hutamatr/GoBlogify/model/web"
	"github.com/julienschmidt/httprouter"
)

type RouterControllers struct {
	UserController     controllersUser.UserController
	PostController     controllersPost.PostController
	CategoryController controllersCategory.CategoryController
	RoleController     controllersRole.RoleController
	CommentController  controllersComment.CommentController
}

func Router(routerControllers *RouterControllers) *httprouter.Router {
	router := httprouter.New()

	router.POST("/api/signup", routerControllers.UserController.CreateUser)
	router.POST("/api/signin", routerControllers.UserController.SignInUser)
	router.POST("/api/signout", routerControllers.UserController.SignOutUser)
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

	router.POST("/api/post", routerControllers.PostController.CreatePost)
	router.GET("/api/post", routerControllers.PostController.FindAllPost)
	router.GET("/api/post/:postId", routerControllers.PostController.FindByIdPost)
	router.PUT("/api/post/:postId", routerControllers.PostController.UpdatePost)
	router.DELETE("/api/post/:postId", routerControllers.PostController.DeletePost)

	router.POST("/api/comment", routerControllers.CommentController.CreateComment)
	router.GET("/api/comment", routerControllers.CommentController.FindCommentsByPost)
	router.GET("/api/comment/:commentId", routerControllers.CommentController.FindByIdComment)
	router.PUT("/api/comment/:commentId", routerControllers.CommentController.UpdateComment)
	router.DELETE("/api/comment/:commentId", routerControllers.CommentController.DeleteComment)

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
