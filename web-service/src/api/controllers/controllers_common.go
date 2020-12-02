package api

import (
	"log"

	config "web-service/src/config"
	containers "web-service/src/storage_container"
	utils "web-service/src/utils"
)

var (
	warningLogger  *log.Logger
	errorLogger    *log.Logger
	debugLogger    *log.Logger
	uploadFilesDir string
)

func InitializeControllers(container containers.ClientContainer) {
	warningLogger = utils.WarningLogger
	errorLogger = utils.ErrorLogger
	debugLogger = utils.DebugLogger

	uploadFilesDir = config.Internal.UploadFilesDir

	initializeViewRoomController(container)
	initializeUploadFilesController()
}
