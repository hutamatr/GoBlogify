package helpers

import (
	"github.com/rs/cors"
)

func Cors() *cors.Cors {
	env := NewEnv()
	appEnv := env.App.AppEnv

	cors := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedOrigins: []string{
			"http://localhost:5173",
			"http://127.0.0.1:5173",
			"http://localhost:4173",
			"http://127.0.0.1:4173",
		},
		AllowedHeaders:   []string{"Origin", "Authorization", "Content-Type"},
		AllowCredentials: true,
		Debug:            appEnv == "development",
	})

	return cors
}
