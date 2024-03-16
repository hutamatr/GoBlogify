package helpers

import (
	"log"
	"os"
	"time"

	"github.com/fatih/color"
)

func ServerRunningText() {

	env := os.Getenv("APP_ENV")
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	textColor := color.New(color.FgGreen).Add(color.BgBlack).SprintfFunc()
	serverRunningText := []string{"Connecting to database...", "Connected!", "Server running on http://" + host + ":" + port, "Environment: " + env}
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for i := range serverRunningText {
		<-ticker.C
		log.Printf("%s %s\n", textColor("OK"), serverRunningText[i])
	}
}
