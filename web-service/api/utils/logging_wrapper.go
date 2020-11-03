package utils

import (
	"path/filepath"
	"os"
	"log"
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

	loggingFile, err = os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	DebugLogger = log.New(loggingFile, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	DebugLogger.Println("Initialized logging to", path)
}

func GetLogger(prefix string) *log.Logger {
	return log.New(loggingFile, prefix, log.Ldate|log.Ltime|log.Lshortfile)
}









