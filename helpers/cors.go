package helpers

import (
	"os"

	"github.com/rs/cors"
)

func Cors() *cors.Cors {
	var AppEnv = os.Getenv("APP_ENV")

	cors := cors.New(cors.Options{
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedOrigins:   []string{"http://localhost:8080", "http://127.0.0.1:8080"},
		AllowedHeaders:   []string{"Origin", "Authorization", "Content-Type"},
		AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in production
		Debug: AppEnv == "development",
	})

	return cors
}
