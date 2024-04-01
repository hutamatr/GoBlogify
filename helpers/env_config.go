package helpers

import (
	"os"

	"github.com/joho/godotenv"
)

type App struct {
	AppEnv string
	Host   string
	Port   string
}

type DB struct {
	Host     string
	Port     string
	Username string
	Password string
	DbName   string
}

type SecretToken struct {
	AccessSecret  string
	RefreshSecret string
}

type Auth struct {
	AdminCode string
}

type Env struct {
	App         *App
	DB          *DB
	SecretToken *SecretToken
	Auth        *Auth
}

func init() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	godotenv.Load(".env." + env)
	if env != "test" {
		godotenv.Load(".env")
	}
	godotenv.Load(".env." + env)
	godotenv.Load()
}

func NewEnv() *Env {
	return &Env{
		App: &App{
			AppEnv: os.Getenv("APP_ENV"),
			Host:   os.Getenv("HOST"),
			Port:   os.Getenv("PORT"),
		},
		DB: &DB{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			Username: os.Getenv("DB_USERNAME"),
			Password: os.Getenv("DB_PASSWORD"),
			DbName:   os.Getenv("DB_NAME"),
		},
		SecretToken: &SecretToken{
			AccessSecret:  os.Getenv("ACCESS_TOKEN_SECRET"),
			RefreshSecret: os.Getenv("REFRESH_TOKEN_SECRET"),
		},
		Auth: &Auth{
			AdminCode: os.Getenv("ADMIN_CODE"),
		},
	}
}
