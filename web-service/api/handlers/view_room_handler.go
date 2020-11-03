package handlers

import (
	"net/http"
	"fmt"
	guuid "github.com/google/uuid"
	nlp   "web-service/api/python_wrappers"
)

var (
	UUIDResult     map[guuid.UUID][][]float32
)

func PrepareViewRoomHandler() {
	UUIDResult = make(map[guuid.UUID][][]float32)
}

func prepareViewForUUID(id guuid.UUID) {
	res, err := nlp.GetPairwiseSimilarity(UploadFilesDir, "--external")
	if err != nil {
		ErrorLogger.Println(err)
	}
	UUIDResult[id] = res
}

func ViewRoomHandler(w http.ResponseWriter, req *http.Request) {
	DebugLogger.Println("viewRoom Endpoint hit")

	// Retrieve view UUID.
	var id_str = req.URL.Query().Get("id")
	if id_str == "" {
		ErrorLogger.Println("Invalid url format")
		http.Error(w, "Incorrect form of url: hostname/view?id=UUID_VALUE expected", http.StatusBadRequest)
		return
	}

	var err error
	view_id, err := guuid.Parse(id_str)
	if err != nil {
		ErrorLogger.Println("Invalid UUID value")
		http.Error(w, "Invalid UUID value", http.StatusBadRequest)
		return
	}

	// If we are not ready, deny of service and halt
	result, present := UUIDResult[view_id]
	if !present {
		fmt.Fprintf(w, "Result for %s is not found.", view_id);
		return
	}

	// Serve the result.
	fmt.Fprintf(w, "%v", result)
}