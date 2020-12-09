package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"

	s3support "web-service/src/s3support"

	guuid "github.com/google/uuid"
	nlp "web-service/src/text_similarity"
)

func CompareFilesHandler(w http.ResponseWriter, req *http.Request) {
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

	id, err := guuid.Parse(urlParsedQuery.Get("id"))
	fileNames := []string{urlParsedQuery.Get("f1"), urlParsedQuery.Get("f2")}

	if err != nil || fileNames[0] == "" || fileNames[1] == "" {
		body["Message"] = `Expected url format: /cmp_files?id=UUID4&f1=str&f2=str`
		if err != nil {
			body["Error"] = err.Error()
		} else {
			body["Error"] = fmt.Sprint("Incorrect form of the input url: ", req.URL.RawQuery)
		}
		compriseMsg(w, body, http.StatusUnprocessableEntity)
		logMsg(warningLogger, body, http.StatusUnprocessableEntity)

	} else if result, err := db.GetResValue(id); reportUnreadyClient(w, id, result, err) {
		return
	}

	var contentBuf bytes.Buffer
	var okErr error = nil
	var fileLen [2]int64
	for idx, name := range fileNames {
		r, err := s3support.DownloadFileByUUID(id, name)
		if err != nil {
			okErr = err
			break
		}

		written, err := io.Copy(&contentBuf, r)
		if err != nil {
			okErr = err
			break
		}

		fileLen[idx] = written
	}

	if okErr != nil {
		body["Error"] = okErr.Error()

		compriseMsg(w, body, http.StatusAccepted)
		logMsg(warningLogger, body, http.StatusAccepted)
	} else {
		option, html := urlParsedQuery.Get("option"), urlParsedQuery.Get("html")
		editcost, timeout := urlParsedQuery.Get("editcost"), urlParsedQuery.Get("timeout")

		res, err := nlp.GetFilesDifference(contentBuf, fileLen, option, html, editcost, timeout)
		if err != nil {
			body["Message"] = fmt.Sprintf("Unable to get content difference for %s, %s", fileNames[0], fileNames[1])
			body["Error"] = err.Error()

			compriseMsg(w, body, http.StatusInternalServerError)
			logMsg(errorLogger, body, http.StatusInternalServerError)

		} else {
			body["Result"] = res

			if html == "false" {
				compriseMsg(w, body, http.StatusOK)
			} else {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				fmt.Fprintf(w, res)
			}
			logMsg(debugLogger, body, http.StatusOK)
		}
	}
}
