package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	guuid "github.com/google/uuid"
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

func reportUnreadyClient(w http.ResponseWriter, id guuid.UUID, result *utils.Pair, err error) bool {
	body := make(map[string]interface{})

	if err != nil {
		body["Error"] = fmt.Sprintf("Result for %s is not found.", id)

		compriseMsg(w, body, http.StatusAccepted)
		logMsg(warningLogger, body, http.StatusAccepted)
	} else if result.First == containers.Error {
		body["Message"] = "Error encountered when analyzing the input, please reupload the files"
		body["Error"] = result.Second

		compriseMsg(w, body, http.StatusAccepted)
		logMsg(errorLogger, body, http.StatusAccepted)

	} else if result.First == containers.Pending {
		body["Error"] = "Analyzing the files haven't been completed yet, try again in several minutes"

		compriseMsg(w, body, http.StatusAccepted)
		logMsg(debugLogger, body, http.StatusAccepted)
	} else {
		return false
	}

	return true
}
