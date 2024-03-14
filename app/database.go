package app

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hutamatr/go-blog-api/helpers"
)

func ConnectDB() *sql.DB {
	var DBName = os.Getenv("DB_NAME")
	var DBUsername = os.Getenv("DB_USERNAME")
	var DBPassword = os.Getenv("DB_PASSWORD")
	var DBPort = os.Getenv("DB_PORT")
	var Host = os.Getenv("HOST")

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", DBUsername, DBPassword, Host, DBPort, DBName))
	helpers.PanicError(err)

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(1 * time.Hour)

	return db
}
