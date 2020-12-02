package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	s3support "web-service/s3support"
	config "web-service/src/config"
	containers "web-service/src/storage_container"
	nlp "web-service/src/text_similarity"

	guuid "github.com/google/uuid"
)

var (
	db containers.ClientContainer
)

func InitializeViewRoomController(container containers.ClientContainer) {
	db = container
}

func prepareViewForUUID(id guuid.UUID) {
	db.SavePendingClient(id, "")
	res, err := nlp.GetPairwiseSimilarity(config.UploadFilesDir)
	if err != nil {
		db.SaveErrorClient(id, err.Error())
		ErrorLogger.Println(err)
	} else {
		db.SaveSuccessClient(id, res)
		files, err := ioutil.ReadDir(path.Join(config.UploadFilesDir, id.String()))
		if err != nil {
			ErrorLogger.Println(err)
		}
		for _, file := range files {
			filePath := path.Join(config.UploadFilesDir, id.String(), file.Name())
			file_, _ := os.Open(filePath)
			err = s3support.StoreFileByUUID(id, file_, file.Name())
			file_.Close()
			os.Remove(filePath)
		}
	}
	os.RemoveAll(path.Join(config.UploadFilesDir, id.String()))
}

func ViewRoomHandler(w http.ResponseWriter, req *http.Request) {
	DebugLogger.Println("viewRoom Endpoint hit")

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
	result, present := db.GetResValue(view_id)
	if !present {
		fmt.Fprintf(w, "Result for %s is not found.", view_id)
	} else if result.First == containers.Error {
		fmt.Fprintln(w, "Error encountered when analyzing the input, please reupload the files.")
	} else if result.First == containers.Pending {
		fmt.Fprintln(w, "Analyzing the files haven't been completed yet, try again in several minutes.")
	} else {
		fmt.Fprintf(w, "%v", result.Second)
	}
}
