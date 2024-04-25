package main

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/hutamatr/GoBlogify/database"
	"github.com/hutamatr/GoBlogify/utils"

	"github.com/hutamatr/GoBlogify/helpers"
	"github.com/hutamatr/GoBlogify/middleware"

	"github.com/hutamatr/GoBlogify/routes"

	"log"
)

func main() {
	db := database.ConnectDB()
	validator := validator.New(validator.WithRequiredStructEnabled())

	roleController := utils.InitializedRoleController(db, validator)
	userController := utils.InitializedUserController(db, validator)
	adminController := utils.InitializedAdminController(db, validator)
	postController := utils.InitializedPostController(db, validator)
	commentController := utils.InitializedCommentController(db, validator)
	categoryController := utils.InitializedCategoryController(db, validator)
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
