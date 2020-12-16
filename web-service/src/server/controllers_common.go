package server

import (
	"fmt"
	"net/http"

	guuid "github.com/google/uuid"
	containers "web-service/src/storage_container"
	utils "web-service/src/utils"
)

func (s *Server) reportUnreadyClient(w http.ResponseWriter, id guuid.UUID, result *utils.Pair, err error) bool {
	body := make(map[string]interface{})

	if err != nil {
		body["Error"] = fmt.Sprintf("Result for %s is not found.", id)

		utils.CompriseMsg(w, body, http.StatusAccepted)
		utils.LogMsg(s.warningLogger, body, http.StatusAccepted)
	} else if result.First == containers.Error {
		body["Message"] = "Error encountered when analyzing the input, please reupload the files"
		body["Error"] = result.Second

		utils.CompriseMsg(w, body, http.StatusAccepted)
		utils.LogMsg(s.errorLogger, body, http.StatusAccepted)

	} else if result.First == containers.Pending {
		body["Error"] = "Analyzing the files haven't been completed yet, try again in several minutes"

		utils.CompriseMsg(w, body, http.StatusAccepted)
		utils.LogMsg(s.debugLogger, body, http.StatusAccepted)
	} else {
		return false
	}

	return true
}
