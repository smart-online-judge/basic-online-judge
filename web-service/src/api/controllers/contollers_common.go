package api

import (
	"log"
	"os"
	config "web-service/src/config"
	utils "web-service/src/utils"
)

var (
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
	DebugLogger   *log.Logger
)

func InitializeControllersCommon() {
	WarningLogger = utils.GetLogger("WARNING: ")
	ErrorLogger = utils.GetLogger("ERROR: ")
	DebugLogger = utils.GetLogger("DEBUG: ")
}

func InitializeUploadFilesController() {
	if _, err := os.Stat(config.UploadFilesDir); os.IsNotExist(err) {
		if err := os.MkdirAll(config.UploadFilesDir, os.ModePerm); err != nil {
			ErrorLogger.Fatal(err)
		}
	}
}
