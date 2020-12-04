package api

import (
	"fmt"
	"net/http"
	"net/url"

	s3support "web-service/src/s3support"

	guuid "github.com/google/uuid"
)

func GetFileLinkById(w http.ResponseWriter, req *http.Request) {
	body := make(map[string]interface{})

	if req.Method != "GET" {
		body["Error"] = fmt.Sprintf("%s request type isn't supported", req.Method)

		compriseMsg(w, body, http.StatusMethodNotAllowed)
		logMsg(warningLogger, body, http.StatusMethodNotAllowed)
		return
	}

	urlParsedQuery, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		body["Message"] = "Unable to parse input url"
		body["Error"] = err.Error()

		compriseMsg(w, body, http.StatusUnprocessableEntity)
		logMsg(warningLogger, body, http.StatusUnprocessableEntity)
		return
	}

	id_str, fileName := urlParsedQuery.Get("id"), urlParsedQuery.Get("name")
	id, err := guuid.Parse(id_str)
	if err != nil {
		body["Message"] = fmt.Sprintf("Invalid id(%s) value: UUID4 expected", id_str)
		body["Error"] = err.Error()

		compriseMsg(w, body, http.StatusUnprocessableEntity)
		logMsg(warningLogger, body, http.StatusUnprocessableEntity)
		return

	} else if !db.ClientExists(id){
		body["Error"] = fmt.Sprintf("Client with given id(%s) wasn't found", id.String())

		compriseMsg(w, body, http.StatusAccepted)
		logMsg(warningLogger, body, http.StatusAccepted)
		return

	} else if fileName == "" {
		body["Error"] = fmt.Sprintf("Invalid filename(%s)", fileName)

		compriseMsg(w, body, http.StatusUnprocessableEntity)
		logMsg(warningLogger, body, http.StatusUnprocessableEntity)
		return
	}

	presignedURL := s3support.PrepareViewFileURL(id, fileName)
	if presignedURL == nil {
		body["Error"] = fmt.Sprintf("Unable to find a file with such name %s", fileName)

		compriseMsg(w, body, http.StatusAccepted)
		logMsg(warningLogger, body, http.StatusAccepted)

	} else {
		body["Link"] = presignedURL.String()
		compriseMsg(w, body, http.StatusOK)
		logMsg(debugLogger, body, http.StatusOK)
	}
}
