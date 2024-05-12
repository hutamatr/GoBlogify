package main

import (
	"net/http"

	"github.com/hutamatr/GoBlogify/database"
	"github.com/hutamatr/GoBlogify/exception"
	"github.com/hutamatr/GoBlogify/utils"

	"github.com/hutamatr/GoBlogify/helpers"
	"github.com/hutamatr/GoBlogify/middleware"

	"github.com/hutamatr/GoBlogify/routes"

	"log"
)

func main() {
	db := database.ConnectDB()
	if err := helpers.CustomValidation(); err != nil {
		panic(exception.NewBadRequestError("invalid request"))
	}

	roleController := utils.InitializedRoleController(db, helpers.Validate)
	userController := utils.InitializedUserController(db, helpers.Validate)
	adminController := utils.InitializedAdminController(db, helpers.Validate)
	postController := utils.InitializedPostController(db, helpers.Validate)
	commentController := utils.InitializedCommentController(db, helpers.Validate)
	categoryController := utils.InitializedCategoryController(db, helpers.Validate)
	followController := utils.InitializedFollowController(db)

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

	if err := http.ListenAndServe(server.Addr, helpers.LogRequest(server.Handler)); err != nil {
		log.Fatal("Server failed to start")
		return
	}
}
