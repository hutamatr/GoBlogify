package go_blog

import (
	"database/sql"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hutamatr/go-blog-api/cmd/go_blog/helper"
)

func ConnectDB() *sql.DB {
	var DBName = os.Getenv("DB_NAME")
	var DBUsername = os.Getenv("DB_USERNAME")
	var DBPassword = os.Getenv("DB_PASSWORD")
	var DBPort = os.Getenv("DB_PORT")
	var Host = os.Getenv("HOST")

	db, err := sql.Open("mysql", DBUsername+":"+DBPassword+"@tcp("+Host+":"+DBPort+")/"+DBName+"?parseTime=true")
	helper.PanicError(err)

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(1 * time.Hour)

	return db
}
