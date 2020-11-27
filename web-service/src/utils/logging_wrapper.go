package utils

import (
	"log"
	"os"
	"path/filepath"
)

var (
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
	DebugLogger   *log.Logger
	loggingFile   *os.File
)

func InitializeLogger(path string) {
	dir, _ := filepath.Split(path)

	var err error
	if err = os.MkdirAll(dir, 0777); err != nil {
		log.Fatal(err)
	}

	loggingFile, err = os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	DebugLogger = log.New(loggingFile, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	DebugLogger.Println("Initialized logging to", path)
}

func GetLogger(prefix string) *log.Logger {
	return log.New(loggingFile, prefix, log.Ldate|log.Ltime|log.Lshortfile)
}

// Get a logger with custom storage path
func GetLoggerPkgScoped(prefix string, storagePath *os.File) *log.Logger {
	return log.New(storagePath, prefix, log.Ldate|log.Ltime|log.Lshortfile)
}
