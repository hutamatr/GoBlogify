package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hutamatr/GoBlogify/helpers"
)

func ConnectDB() *sql.DB {
	env := helpers.NewEnv()
	DBName := env.DB.DbName
	DBUsername := env.DB.Username
	DBPassword := env.DB.Password
	DBPort := env.DB.Port
	Host := env.DB.Host

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", DBUsername, DBPassword, Host, DBPort, DBName))
	helpers.PanicError(err, "failed to connect database")

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(1 * time.Hour)

	return db
}
