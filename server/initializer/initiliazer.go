package initializer

import (
	"log"
	"os"
	"time"
)

func InitializeLogger() *log.Logger {
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}

	return log.New(logFile, "APP: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func LogMessage(logger *log.Logger, functionName, message string) {
	logger.Printf("%s %s: %s\n", time.Now().Format("2006/01/02 15:04:05"), functionName, message)
}
