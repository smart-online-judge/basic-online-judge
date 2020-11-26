package handlers

import (
	"fmt"
	guuid "github.com/google/uuid"
	"net/http"
	nlp "web-service/api/text_similarity"
	containers "web-service/api/storage_container"
	"strings"
)

var (
	db containers.ClientContainer
)

func InitializeViewRoomHandler(container containers.ClientContainer) {
	db = container
}

func prepareViewForUUID(id guuid.UUID) {
	res, err := nlp.GetPairwiseSimilarity(UploadFilesDir)
	if err != nil {
		ErrorLogger.Println(err)
	}

	// TODO: Here we may need to somehow handle the err status
	db.SaveClient(id, res)
}

func ViewRoomHandler(w http.ResponseWriter, req *http.Request) {
	DebugLogger.Println("viewRoom Endpoint hit")

	// Retrieve view id.
	id_str := strings.TrimPrefix(req.URL.Path, "/view/")
	if (id_str == "" || strings.Contains(id_str, "/")) {
		http.Error(w, "Incorrect form of url: hostname/view/{id} expected", http.StatusBadRequest)
		return
	}

	view_id, err := guuid.Parse(id_str)
	if err != nil {
		http.Error(w, "Invalid id value: UUID4 expected", http.StatusBadRequest)
		return
	}

	// If we are not ready, deny of service and halt
	result, present := db.GetValue(view_id)
	if !present {
		fmt.Fprintf(w, "Result for %s is not found.", view_id)
		return
	}

	// Serve the result.
	fmt.Fprintf(w, "%v", result)
}
