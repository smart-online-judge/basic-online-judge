package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	guuid "github.com/google/uuid"
	containers "web-service/src/storage_container"
	utils "web-service/src/utils"
)

func (s *Server) compriseMsg(w http.ResponseWriter, body map[string]interface{}, status int) error {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(body)
	return err
}

func (s *Server) logMsg(logger *log.Logger, body map[string]interface{}, status int) {
	if status != http.StatusOK {
		logger.Println(body["Error"])
	} else {
		logger.Println(body["Result"])
	}
}

func (s *Server) reportUnreadyClient(w http.ResponseWriter, id guuid.UUID, result *utils.Pair, err error) bool {
	body := make(map[string]interface{})

	if err != nil {
		body["Error"] = fmt.Sprintf("Result for %s is not found.", id)

		s.compriseMsg(w, body, http.StatusAccepted)
		s.logMsg(s.warningLogger, body, http.StatusAccepted)
	} else if result.First == containers.Error {
		body["Message"] = "Error encountered when analyzing the input, please reupload the files"
		body["Error"] = result.Second

		s.compriseMsg(w, body, http.StatusAccepted)
		s.logMsg(s.errorLogger, body, http.StatusAccepted)

	} else if result.First == containers.Pending {
		body["Error"] = "Analyzing the files haven't been completed yet, try again in several minutes"

		s.compriseMsg(w, body, http.StatusAccepted)
		s.logMsg(s.debugLogger, body, http.StatusAccepted)
	} else {
		return false
	}

	return true
}
