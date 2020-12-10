package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	guuid "github.com/google/uuid"
)

func (s *Server) prepareViewForUUID(id guuid.UUID) {
	s.db.SavePendingClient(id, "")

	var clientDir = path.Join(s.conf.Internal.UploadFilesDir, id.String())
	res, err := s.nlpCore.GetPairwiseSimilarity(clientDir)
	if err != nil {
		s.db.SaveErrorClient(id, err.Error())
		s.errorLogger.Println(err)
	} else {
		s.db.SaveSuccessClient(id, res)
		files, err := ioutil.ReadDir(clientDir)
		if err != nil {
			s.errorLogger.Println(err)
		} else {
			for _, file := range files {
				s.s3Client.UploadFsFileByUUID(id, clientDir, file.Name())
			}
		}
	}
	os.RemoveAll(clientDir)
}

func (s *Server) viewRoomHandler(w http.ResponseWriter, req *http.Request) {
	s.debugLogger.Println("viewRoom Endpoint hit")
	body := make(map[string]interface{})

	if req.Method != "GET" {
		body["Error"] = fmt.Sprintf("%s request type isn't supported", req.Method)

		s.compriseMsg(w, body, http.StatusMethodNotAllowed)
		s.logMsg(s.warningLogger, body, http.StatusMethodNotAllowed)
		return
	}
	// Retrieve view id.
	id_str := strings.TrimPrefix(req.URL.Path, "/api/view/")
	if id_str == "" || strings.Contains(id_str, "/") {
		body["Error"] = fmt.Sprint("Incorrect form of url", req.URL.Path)
		body["Message"] = "hostname/api/view/{id} expected"

		s.compriseMsg(w, body, http.StatusMethodNotAllowed)
		s.logMsg(s.warningLogger, body, http.StatusMethodNotAllowed)
		return
	}

	view_id, err := guuid.Parse(id_str)
	if err != nil {
		body["Message"] = fmt.Sprintf("Invalid id(%s) value: UUID4 expected", id_str)
		body["Error"] = err.Error()

		s.compriseMsg(w, body, http.StatusUnprocessableEntity)
		s.logMsg(s.warningLogger, body, http.StatusUnprocessableEntity)
		return
	}

	result, err := s.db.GetResValue(view_id)
	if s.reportUnreadyClient(w, view_id, result, err) {
		return

	} else {
		body["Result"] = result.Second
		body["Message"] = "Text similarity matrix"
		body["Files"] = s.s3Client.ListFilesByUUID(view_id)

		s.compriseMsg(w, body, http.StatusOK)
		s.logMsg(s.debugLogger, body, http.StatusOK)
	}
}
