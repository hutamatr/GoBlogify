package helpers

import (
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

var (
	Logger = logrus.New()
)

func init() {
	Logger.SetFormatter(&logrus.JSONFormatter{})
	// Logger.SetLevel(logrus.DebugLevel)
}

func logToFile(filename string) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		Logger.SetOutput(file)
	} else {
		Logger.Info("Failed to log to file, using default stdout")
		Logger.SetOutput(os.Stdout)
	}
}

func LogRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Logger.Infof("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func LogError(format string, a ...interface{}) {
	logToFile("error-log.log")
	Logger.Errorf(format, a...)
}
