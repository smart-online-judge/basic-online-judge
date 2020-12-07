package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	s3support "web-service/src/s3support"
	nlp "web-service/src/text_similarity"

	guuid "github.com/google/uuid"
)

func prepareViewForUUID(id guuid.UUID) {
	db.SavePendingClient(id, "")

	var clientDir = path.Join(uploadFilesDir, id.String())
	res, err := nlp.GetPairwiseSimilarity(clientDir)
	if err != nil {
		db.SaveErrorClient(id, err.Error())
		errorLogger.Println(err)
	} else {
		db.SaveSuccessClient(id, res)
		files, err := ioutil.ReadDir(clientDir)
		if err != nil {
			errorLogger.Println(err)
		} else {
			for _, file := range files {
				s3support.UploadFsFileByUUID(id, clientDir, file.Name())
			}
		}
	}
	os.RemoveAll(clientDir)
}

func ViewRoomHandler(w http.ResponseWriter, req *http.Request) {
	debugLogger.Println("viewRoom Endpoint hit")
	body := make(map[string]interface{})

	if req.Method != "GET" {
		body["Error"] = fmt.Sprintf("%s request type isn't supported", req.Method)

		compriseMsg(w, body, http.StatusMethodNotAllowed)
		logMsg(warningLogger, body, http.StatusMethodNotAllowed)
		return
	}
	// Retrieve view id.
	id_str := strings.TrimPrefix(req.URL.Path, "/api/view/")
	if id_str == "" || strings.Contains(id_str, "/") {
		body["Error"] = fmt.Sprint("Incorrect form of url", req.URL.Path)
		body["Message"] = "hostname/api/view/{id} expected"

		compriseMsg(w, body, http.StatusMethodNotAllowed)
		logMsg(warningLogger, body, http.StatusMethodNotAllowed)
		return
	}

	view_id, err := guuid.Parse(id_str)
	if err != nil {
		body["Message"] = fmt.Sprintf("Invalid id(%s) value: UUID4 expected", id_str)
		body["Error"] = err.Error()

		compriseMsg(w, body, http.StatusUnprocessableEntity)
		logMsg(warningLogger, body, http.StatusUnprocessableEntity)
		return
	}

	result, err := db.GetResValue(view_id)
	if reportUnreadyClient(w, view_id, result, err) {
		return

	} else {
		body["Result"] = result.Second
		body["Message"] = "Text similarity matrix"
		body["Files"] = s3support.ListFilesByUUID(view_id)

		compriseMsg(w, body, http.StatusOK)
		logMsg(debugLogger, body, http.StatusOK)
	}
}
