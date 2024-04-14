package main

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/hutamatr/GoBlogify/admin"
	"github.com/hutamatr/GoBlogify/app"
	"github.com/hutamatr/GoBlogify/category"
	"github.com/hutamatr/GoBlogify/comment"
	"github.com/hutamatr/GoBlogify/follow"
	"github.com/hutamatr/GoBlogify/post"
	"github.com/hutamatr/GoBlogify/role"
	"github.com/hutamatr/GoBlogify/user"

	"github.com/hutamatr/GoBlogify/helpers"
	"github.com/hutamatr/GoBlogify/middleware"

	"github.com/hutamatr/GoBlogify/routes"

	"github.com/rs/zerolog/log"
)

func main() {
	db := app.ConnectDB()
	validator := validator.New(validator.WithRequiredStructEnabled())

	roleRepository := role.NewRoleRepository()
	roleService := role.NewRoleService(roleRepository, db, validator)
	roleController := role.NewRoleController(roleService)

	userRepository := user.NewUserRepository()
	userService := user.NewUserService(userRepository, roleRepository, db, validator)
	userController := user.NewUserController(userService)

	adminService := admin.NewAdminService(userRepository, roleRepository, db, validator)
	adminController := admin.NewAdminControllerImpl(adminService)

	postRepository := post.NewPostRepository()
	postService := post.NewPostService(postRepository, db, validator)
	postController := post.NewPostController(postService)

	commentRepository := comment.NewCommentRepository()
	commentService := comment.NewCommentService(commentRepository, db, validator)
	commentController := comment.NewCommentController(commentService)

	categoryRepository := category.NewCategoryRepository()
	categoryService := category.NewCategoryService(categoryRepository, db, validator)
	categoryController := category.NewCategoryController(categoryService)

	followRepository := follow.NewFollowRepository()
	followService := follow.NewFollowService(followRepository, db)
	followController := follow.NewFollowController(followService)

	router := routes.Router(&routes.RouterControllers{
		Admin:    adminController,
		User:     userController,
		Post:     postController,
		Category: categoryController,
		Role:     roleController,
		Comment:  commentController,
		Follow:   followController,
	})

	cors := helpers.Cors()
	corsHandler := cors.Handler(router)

	server := http.Server{
		Addr:    ":8080",
		Handler: middleware.NewAuthMiddleware(corsHandler),
	}

	helpers.ServerRunningText()

	if err := http.ListenAndServe(server.Addr, server.Handler); err != nil {
		log.Fatal().Err(err).Msg("Server failed to start")
		return
	}
}
