package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	guuid "github.com/google/uuid"
	s3support "web-service/src/s3support"
	containers "web-service/src/storage_container"
	nlp "web-service/src/text_similarity"
)

var (
	db containers.ClientContainer
)

func initializeViewRoomController(container containers.ClientContainer) {
	db = container
}

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

	// Retrieve view id.
	id_str := strings.TrimPrefix(req.URL.Path, "/view/")
	if id_str == "" || strings.Contains(id_str, "/") {
		http.Error(w, "Incorrect form of url: hostname/view/{id} expected", http.StatusBadRequest)
		return
	}

	view_id, err := guuid.Parse(id_str)
	if err != nil {
		http.Error(w, "Invalid id value: UUID4 expected", http.StatusBadRequest)
		return
	}

	// If we are not ready, deny of service and halt
	result, err := db.GetResValue(view_id)
	if err != nil {
		fmt.Fprintf(w, "Result for %s is not found.", view_id)
	} else if result.First == containers.Error {
		fmt.Fprintln(w, "Error encountered when analyzing the input, please reupload the files.")
	} else if result.First == containers.Pending {
		fmt.Fprintln(w, "Analyzing the files haven't been completed yet, try again in several minutes.")
	} else {
		fmt.Fprintf(w, "%v", result.Second)
	}
}
