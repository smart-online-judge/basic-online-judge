package utils

import (
	"log"
	"os"
	"path"

	config "web-service/src/config"
)

const (
	defaultLoggingFileName = "main_log.log"
)

var (
	ErrorLogger        *log.Logger
	WarningLogger      *log.Logger
	DebugLogger        *log.Logger
	defaultloggingFile *os.File
)

func InitializeLogger() error {
	dir := config.Internal.LoggingDir

	var err error
	if err = os.MkdirAll(dir, 0777); err != nil {
		return err
	}

	file, err := os.OpenFile(
		path.Join(dir, defaultLoggingFileName),
		os.O_CREATE|os.O_TRUNC|os.O_WRONLY,
		0666)
	if err != nil {
		return err
	}

	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	DebugLogger = log.New(file, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	return nil
}

// Get a logger with custom storage path
func GetLoggerPkgScoped(prefix string, storagePath *os.File) *log.Logger {
	return log.New(storagePath, prefix, log.Ldate|log.Ltime|log.Lshortfile)
}
