package helpers

import (
	"fmt"
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
		fmt.Printf("%s ", textColor("OK"))
		fmt.Println(serverRunningText[i])
	}
}
