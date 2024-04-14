package helpers

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
)

func ServerRunningText() {
	env := NewEnv()
	appEnv := env.App.AppEnv
	host := env.App.Host
	port := env.App.Port

	serverRunningText := []string{"Running server...", fmt.Sprintf("Server running on http://%s:%s", host, port), fmt.Sprintf("Environment -> %s", appEnv)}
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for i := range serverRunningText {
		<-ticker.C
		log.Info().Msg(serverRunningText[i])
	}
}
