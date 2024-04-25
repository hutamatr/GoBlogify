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

	router.POST("/api/v1/signup-admin", route.Admin.CreateAdminHandler)
	router.POST("/api/v1/signin-admin", route.Admin.SignInAdminHandler)

	router.POST("/api/v1/signup", route.User.CreateUserHandler)
	router.POST("/api/v1/signin", route.User.SignInUserHandler)
	router.POST("/api/v1/signout", route.User.SignOutUserHandler)
	router.GET("/api/v1/refresh", route.User.GetRefreshTokenHandler)

	router.GET("/api/v1/users", route.User.FindAllUserHandler)
	router.GET("/api/v1/users/:userId", route.User.FindByIdUserHandler)
	router.PUT("/api/v1/users/:userId", route.User.UpdateUserHandler)
	router.DELETE("/api/v1/users/:userId", route.User.DeleteUserHandler)

	router.POST("/api/v1/users/:userId/follow/:toUserId", route.Follow.FollowUserHandler)
	router.DELETE("/api/v1/users/:userId/unfollow/:toUserId", route.Follow.UnfollowUserHandler)
	router.GET("/api/v1/users/:userId/follower", route.Follow.FindAllFollowerByUserHandler)
	router.GET("/api/v1/users/:userId/following", route.Follow.FindAllFollowedByUserHandler)

	router.POST("/api/v1/roles", route.Role.CreateRoleHandler)
	router.GET("/api/v1/roles", route.Role.FindAllRoleHandler)
	router.GET("/api/v1/roles/:roleId", route.Role.FindRoleByIdHandler)
	router.PUT("/api/v1/roles/:roleId", route.Role.UpdateRoleHandler)
	router.DELETE("/api/v1/roles/:roleId", route.Role.DeleteRoleHandler)

	router.POST("/api/v1/posts", route.Post.CreatePostHandler)
	router.GET("/api/v1/posts/:userId", route.Post.FindAllPostByUserHandler)
	router.GET("/api/v1/posts/:userId/following", route.Post.FindAllPostByFollowedHandler)
	router.GET("/api/v1/post/:postId", route.Post.FindByIdPostHandler)
	router.PUT("/api/v1/posts/:postId", route.Post.UpdatePostHandler)
	router.DELETE("/api/v1/posts/:postId", route.Post.DeletePostHandler)

	router.POST("/api/v1/comments", route.Comment.CreateCommentHandler)
	router.GET("/api/v1/comments", route.Comment.FindCommentsByPostHandler)
	router.GET("/api/v1/comments/:commentId", route.Comment.FindCommentByIdHandler)
	router.PUT("/api/v1/comments/:commentId", route.Comment.UpdateCommentHandler)
	router.DELETE("/api/v1/comments/:commentId", route.Comment.DeleteCommentHandler)

	router.POST("/api/v1/categories", route.Category.CreateCategoryHandler)
	router.GET("/api/v1/categories", route.Category.FindAllCategoryHandler)
	router.GET("/api/v1/categories/:categoryId", route.Category.FindByIdCategoryHandler)
	router.PUT("/api/v1/categories/:categoryId", route.Category.UpdateCategoryHandler)
	router.DELETE("/api/v1/categories/:categoryId", route.Category.DeleteCategoryHandler)

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
