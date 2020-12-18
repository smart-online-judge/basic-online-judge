package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	guuid "github.com/google/uuid"
	utils "web-service/src/utils"
)

func (s *Server) compareFilesHandler(w http.ResponseWriter, req *http.Request) {
	body := make(map[string]interface{})

	if req.Method != "GET" {
		body["Error"] = fmt.Sprintf("%s request type isn't supported", req.Method)

		utils.CompriseMsg(w, body, http.StatusMethodNotAllowed)
		utils.LogMsg(s.warningLogger, body, http.StatusMethodNotAllowed)
		return
	}

	urlParsedQuery, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		body["Message"] = "Unable to parse input url"
		body["Error"] = err.Error()

		utils.CompriseMsg(w, body, http.StatusUnprocessableEntity)
		utils.LogMsg(s.warningLogger, body, http.StatusUnprocessableEntity)
		return
	}

	id, err := guuid.Parse(urlParsedQuery.Get("id"))
	fileNames := []string{urlParsedQuery.Get("f1"), urlParsedQuery.Get("f2")}
	if err != nil || fileNames[0] == "" || fileNames[1] == "" {
		body["Message"] = `Expected url format: /cmp_files?id=UUID4&f1=str&f2=str`
		if err != nil {
			body["Error"] = err.Error()
		} else {
			body["Error"] = fmt.Sprint("Incorrect form of the input url: ", req.URL.String())
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		enc := json.NewEncoder(w)
		enc.SetEscapeHTML(false)
		enc.Encode(body)

		utils.LogMsg(s.warningLogger, body, http.StatusUnprocessableEntity)
		return

	} else if result, err := s.db.GetResValue(id); s.reportUnreadyClient(w, id, result, err) {
		return
	}

	var contentBuf bytes.Buffer
	var okErr error = nil
	var fileLen [2]int64
	for idx, name := range fileNames {
		r, err := s.s3Client.DownloadFileByUUID(id, name)
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

		utils.CompriseMsg(w, body, http.StatusAccepted)
		utils.LogMsg(s.warningLogger, body, http.StatusAccepted)
	} else {
		option, html := urlParsedQuery.Get("option"), urlParsedQuery.Get("html")
		editcost, timeout := urlParsedQuery.Get("editcost"), urlParsedQuery.Get("timeout")

		res, err := s.nlpCore.GetFilesDifference(contentBuf, fileLen, option, html, editcost, timeout)
		if err != nil {
			body["Message"] = fmt.Sprintf("Unable to get content difference for %s, %s", fileNames[0], fileNames[1])
			body["Error"] = err.Error()

			utils.CompriseMsg(w, body, http.StatusInternalServerError)
			utils.LogMsg(s.errorLogger, body, http.StatusInternalServerError)

		} else {
			body["Result"] = res

			if html == "false" {
				utils.CompriseMsg(w, body, http.StatusOK)
			} else {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				fmt.Fprintf(w, res)
			}
			utils.LogMsg(s.debugLogger, body, http.StatusOK)
		}
	}
}
