package test

import (
	"database/sql"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hutamatr/GoBlogify/routes"
	"github.com/hutamatr/GoBlogify/utils"

	"github.com/hutamatr/GoBlogify/helpers"
	"github.com/hutamatr/GoBlogify/middleware"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load("../.env.test")
	helpers.PanicError(err, "failed to load .env.test")
}

func ConnectDBTest() *sql.DB {
	env := helpers.NewEnv()
	DBName := env.DB.DbName
	DBUsername := env.DB.Username
	DBPassword := env.DB.Password
	DBPort := env.DB.Port
	Host := env.DB.Host

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", DBUsername, DBPassword, Host, DBPort, DBName))
	helpers.PanicError(err, "failed to connect test database")

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(1 * time.Hour)

	return db
}

func DeleteDBTest(db *sql.DB) {
	_, err := db.Exec("DELETE FROM comments")
	helpers.PanicError(err, "failed to delete comment")
	_, err = db.Exec("DELETE FROM post_images")
	helpers.PanicError(err, "failed to delete post_images")
	_, err = db.Exec("DELETE FROM posts")
	helpers.PanicError(err, "failed to delete posts")
	_, err = db.Exec("DELETE FROM categories")
	helpers.PanicError(err, "failed to delete categories")
	_, err = db.Exec("DELETE FROM follows")
	helpers.PanicError(err, "failed to delete follows")
	_, err = db.Exec("DELETE FROM users")
	helpers.PanicError(err, "failed to delete users")
	_, err = db.Exec("DELETE FROM roles")
	helpers.PanicError(err, "failed to delete roles")
}

func SetupRouterTest(db *sql.DB) http.Handler {
	helpers.CustomValidation()

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

	return middleware.NewAuthMiddleware(router)
}

func AddFile(t *testing.T, w *multipart.Writer, fieldName, fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		t.Fatalf("Unable to open file %s: %v", fileName, err)
	}
	defer file.Close()

	fw, err := w.CreateFormFile(fieldName, fileName)
	if err != nil {
		t.Fatalf("Error creating form file: %v", err)
	}

	_, err = io.Copy(fw, file)
	if err != nil {
		t.Fatalf("Error copying file: %v", err)
	}
}
