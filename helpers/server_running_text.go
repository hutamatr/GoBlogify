package helpers

import (
	"fmt"
	"log"
	"time"

	"github.com/fatih/color"
)

func ServerRunningText() {
	env := NewEnv()
	appEnv := env.App.AppEnv
	host := env.App.Host
	port := env.App.Port

	textColor := color.New(color.FgGreen).Add(color.BgBlack).SprintfFunc()
	serverRunningText := []string{"Connecting to database...", "Connected!", fmt.Sprintf("Server running on http://%s:%s", host, port), fmt.Sprintf("Environment: %s", appEnv)}
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for i := range serverRunningText {
		<-ticker.C
		log.Printf("%s %s\n", textColor("OK"), serverRunningText[i])
	}
}
