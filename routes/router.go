package routes

import (
	"net/http"

	"github.com/hutamatr/GoBlogify/admin"
	"github.com/hutamatr/GoBlogify/category"
	"github.com/hutamatr/GoBlogify/comment"
	"github.com/hutamatr/GoBlogify/exception"
	"github.com/hutamatr/GoBlogify/follow"
	"github.com/hutamatr/GoBlogify/helpers"
	"github.com/hutamatr/GoBlogify/post"
	"github.com/hutamatr/GoBlogify/role"
	"github.com/hutamatr/GoBlogify/user"
	"github.com/julienschmidt/httprouter"
)

type RouterControllers struct {
	Admin    admin.AdminController
	User     user.UserController
	Post     post.PostController
	Category category.CategoryController
	Role     role.RoleController
	Comment  comment.CommentController
	Follow   follow.FollowController
}

func Router(route *RouterControllers) *httprouter.Router {
	router := httprouter.New()

	router.POST("/api/signup-admin", route.Admin.CreateAdminHandler)
	router.POST("/api/signin-admin", route.Admin.SignInAdminHandler)

	router.POST("/api/signup", route.User.CreateUserHandler)
	router.POST("/api/signin", route.User.SignInUserHandler)
	router.POST("/api/signout", route.User.SignOutUserHandler)
	router.GET("/api/refresh", route.User.GetRefreshTokenHandler)

	router.GET("/api/user", route.User.FindAllUserHandler)
	router.GET("/api/user/:userId", route.User.FindByIdUserHandler)
	router.PUT("/api/user/:userId", route.User.UpdateUserHandler)
	router.DELETE("/api/user/:userId", route.User.DeleteUserHandler)

	router.POST("/api/user/:userId/follow/:toUserId", route.Follow.FollowUserHandler)
	router.DELETE("/api/user/:userId/unfollow/:toUserId", route.Follow.UnfollowUserHandler)
	router.GET("/api/user/:userId/follower", route.Follow.FindAllFollowerByUserHandler)
	router.GET("/api/user/:userId/following", route.Follow.FindAllFollowedByUserHandler)

	router.POST("/api/role", route.Role.CreateRoleHandler)
	router.GET("/api/role", route.Role.FindAllRoleHandler)
	router.GET("/api/role/:roleId", route.Role.FindRoleByIdHandler)
	router.PUT("/api/role/:roleId", route.Role.UpdateRoleHandler)
	router.DELETE("/api/role/:roleId", route.Role.DeleteRoleHandler)

	router.POST("/api/post", route.Post.CreatePostHandler)
	router.GET("/api/post", route.Post.FindAllPostHandler)
	router.GET("/api/post/:postId", route.Post.FindByIdPostHandler)
	router.PUT("/api/post/:postId", route.Post.UpdatePostHandler)
	router.DELETE("/api/post/:postId", route.Post.DeletePostHandler)

	router.POST("/api/comment", route.Comment.CreateCommentHandler)
	router.GET("/api/comment", route.Comment.FindCommentsByPostHandler)
	router.GET("/api/comment/:commentId", route.Comment.FindCommentByIdHandler)
	router.PUT("/api/comment/:commentId", route.Comment.UpdateCommentHandler)
	router.DELETE("/api/comment/:commentId", route.Comment.DeleteCommentHandler)

	router.POST("/api/category", route.Category.CreateCategoryHandler)
	router.GET("/api/category", route.Category.FindAllCategoryHandler)
	router.GET("/api/category/:categoryId", route.Category.FindByIdCategoryHandler)
	router.PUT("/api/category/:categoryId", route.Category.UpdateCategoryHandler)
	router.DELETE("/api/category/:categoryId", route.Category.DeleteCategoryHandler)

	router.NotFound = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusNotFound)
		userResponse := helpers.ResponseJSON{
			Code:   http.StatusNotFound,
			Status: "NOT FOUND",
		}

		helpers.EncodeJSONFromResponse(writer, userResponse)
	})

	router.PanicHandler = exception.ErrorHandler

	return router
}
