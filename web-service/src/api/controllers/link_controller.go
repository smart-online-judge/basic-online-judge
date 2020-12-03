package api

import (
	"fmt"
	"net/http"
	"net/url"

	s3support "web-service/src/s3support"

	guuid "github.com/google/uuid"
)

/*
   type JSONResponse struct {
	name string
        }
*/

func ServeFileViewByUUID(w http.ResponseWriter, req *http.Request) {
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
		errorLogger.Println(err)
		msg := fmt.Sprintf("Unable to find a file with such name %s", fileName)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	jsonEncoded := fmt.Sprintf("{\"link\": \"%s\"}", presignedURL.String())
	w.Write([]byte(jsonEncoded))
}
