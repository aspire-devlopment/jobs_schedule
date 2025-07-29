package internal

import (
	"log"
	"os"
)

// errorLog is the file handle used for logging errors to a file
var errorLog *os.File

func InitLogger(logPath string) {
	var err error
	errorLog, err = os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	// Set the default logger's output to the errorLog file
	log.SetOutput(errorLog)
}

// LogError writes an error message prefixed with "[ERROR]" to the configured log output.
func LogError(msg string) {
	log.Println("[ERROR] " + msg)
}
