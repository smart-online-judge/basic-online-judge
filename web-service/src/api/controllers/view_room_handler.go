package controllers

import (
	"fmt"
	guuid "github.com/google/uuid"
	"net/http"
	"strings"
	containers "web-service/src/storage_container"
	nlp "web-service/src/text_similarity"
)

var (
	db containers.ClientContainer
)

func InitializeViewRoomHandler(container containers.ClientContainer) {
	db = container
}

func prepareViewForUUID(id guuid.UUID) {
	db.SavePendingClient(id, "Analyzing the input files")

	res, err := nlp.GetPairwiseSimilarity(UploadFilesDir)
	if err != nil {
		db.SaveErrorClient(id, err.Error())
		ErrorLogger.Println(err)
	} else {
		db.SaveResClient(id, res)
	}
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
