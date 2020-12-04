package api

import (
	"fmt"
	"net/http"
	"net/url"

	s3support "web-service/src/s3support"

	guuid "github.com/google/uuid"
)

func GetFileLinkById(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		msg := "Only accepts GET"
		debugLogger.Println(msg)
		http.Error(w, msg, http.StatusMethodNotAllowed)
		return
	}
	urlParsedQuery, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		errorLogger.Println(err)
		http.Error(w, "Unable to parse input files", http.StatusUnprocessableEntity)
		return
	}
	id_str := urlParsedQuery.Get("id")
	fileName := urlParsedQuery.Get("name")
	id, err := guuid.Parse(id_str)
	if err != nil {
		errorLogger.Println(err)
		http.Error(w, "Please provide a valid UUID4", http.StatusUnprocessableEntity)
		return
	}

	presignedURL := s3support.PrepareViewFileURL(id, fileName)
	if presignedURL == nil {
		msg := fmt.Sprintf("Unable to find a file with such name %s", fileName)
		errorLogger.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsonEncoded := fmt.Sprintf("{\"link\": \"%s\"}", presignedURL.String())
	w.Write([]byte(jsonEncoded))
}
