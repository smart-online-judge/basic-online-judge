package api

import (
	"encoding/json"
	"log"
	"net/http"

	config "web-service/src/config"
	containers "web-service/src/storage_container"
	utils "web-service/src/utils"
)

var (
	warningLogger  *log.Logger
	errorLogger    *log.Logger
	debugLogger    *log.Logger
	uploadFilesDir string
	db             containers.ClientContainer
)

func InitializeControllers(container containers.ClientContainer) {
	warningLogger = utils.WarningLogger
	errorLogger = utils.ErrorLogger
	debugLogger = utils.DebugLogger

	uploadFilesDir = config.Internal.UploadFilesDir
	db = container

	initializeUploadFilesController()
}

func compriseMsg(w http.ResponseWriter, body map[string]interface{}, status int) error {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(body)
	return err
}

func logMsg(logger *log.Logger, body map[string]interface{}, status int) {
	if status != http.StatusOK {
		logger.Println(body["Error"])
	} else {
		logger.Println(body["Result"])
	}
}
