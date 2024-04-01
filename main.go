package main

import (
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/hutamatr/GoBlogify/app"
	controllersAdmin "github.com/hutamatr/GoBlogify/controllers/admin"
	controllersCategory "github.com/hutamatr/GoBlogify/controllers/category"
	controllersComment "github.com/hutamatr/GoBlogify/controllers/comment"
	controllersPost "github.com/hutamatr/GoBlogify/controllers/post"
	controllersRole "github.com/hutamatr/GoBlogify/controllers/role"
	controllersUser "github.com/hutamatr/GoBlogify/controllers/user"
	"github.com/hutamatr/GoBlogify/helpers"
	"github.com/hutamatr/GoBlogify/middleware"
	repositoriesCategory "github.com/hutamatr/GoBlogify/repositories/category"
	repositoriesComment "github.com/hutamatr/GoBlogify/repositories/comment"
	repositoriesPost "github.com/hutamatr/GoBlogify/repositories/post"
	repositoriesRole "github.com/hutamatr/GoBlogify/repositories/role"
	repositoriesUser "github.com/hutamatr/GoBlogify/repositories/user"
	"github.com/hutamatr/GoBlogify/routes"
	servicesAdmin "github.com/hutamatr/GoBlogify/services/admin"
	servicesCategory "github.com/hutamatr/GoBlogify/services/category"
	servicesComment "github.com/hutamatr/GoBlogify/services/comment"
	servicesPost "github.com/hutamatr/GoBlogify/services/post"
	servicesRole "github.com/hutamatr/GoBlogify/services/role"
	servicesUser "github.com/hutamatr/GoBlogify/services/user"
)

func main() {
	db := app.ConnectDB()
	validator := validator.New(validator.WithRequiredStructEnabled())

	roleRepository := repositoriesRole.NewRoleRepository()
	roleService := servicesRole.NewRoleService(roleRepository, db, validator)
	roleController := controllersRole.NewRoleController(roleService)

	adminRepository := repositoriesUser.NewUserRepository()
	adminService := servicesAdmin.NewAdminService(adminRepository, roleRepository, db, validator)
	adminController := controllersAdmin.NewAdminControllerImpl(adminService)

	userRepository := repositoriesUser.NewUserRepository()
	userService := servicesUser.NewUserService(userRepository, roleRepository, db, validator)
	userController := controllersUser.NewUserController(userService)

	postRepository := repositoriesPost.NewPostRepository()
	postService := servicesPost.NewPostService(postRepository, db, validator)
	postController := controllersPost.NewPostController(postService)

	commentRepository := repositoriesComment.NewCommentRepository()
	commentService := servicesComment.NewCommentService(commentRepository, db, validator)
	commentController := controllersComment.NewCommentController(commentService)

	categoryRepository := repositoriesCategory.NewCategoryRepository()
	categoryService := servicesCategory.NewCategoryService(categoryRepository, db, validator)
	categoryController := controllersCategory.NewCategoryController(categoryService)

	router := routes.Router(&routes.RouterControllers{
		UserController:     userController,
		PostController:     postController,
		CategoryController: categoryController,
		RoleController:     roleController,
		CommentController:  commentController,
		AdminController:    adminController,
	})

	cors := helpers.Cors()
	corsHandler := cors.Handler(router)

	server := http.Server{
		Addr:    ":8080",
		Handler: middleware.NewAuthMiddleware(corsHandler),
	}

	helpers.ServerRunningText()

	log.Fatal(server.ListenAndServe())
}
