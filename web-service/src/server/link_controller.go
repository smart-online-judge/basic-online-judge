package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	guuid "github.com/google/uuid"
)

func (s *Server) getFileLinkById(w http.ResponseWriter, req *http.Request) {
	body := make(map[string]interface{})

	if req.Method != "GET" {
		body["Error"] = fmt.Sprintf("%s request type isn't supported", req.Method)

		s.compriseMsg(w, body, http.StatusMethodNotAllowed)
		s.logMsg(s.warningLogger, body, http.StatusMethodNotAllowed)
		return
	}

	urlParsedQuery, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		body["Message"] = "Unable to parse input url"
		body["Error"] = err.Error()

		s.compriseMsg(w, body, http.StatusUnprocessableEntity)
		s.logMsg(s.warningLogger, body, http.StatusUnprocessableEntity)
		return
	}

	id_str, fileName := urlParsedQuery.Get("id"), urlParsedQuery.Get("name")
	id, err := guuid.Parse(id_str)

	if err != nil {
		body["Message"] = fmt.Sprintf("Invalid id(%s) value: UUID4 expected", id_str)
		body["Error"] = err.Error()

		s.compriseMsg(w, body, http.StatusUnprocessableEntity)
		s.logMsg(s.warningLogger, body, http.StatusUnprocessableEntity)
		return

	} else if fileName == "" {
		body["Error"] = fmt.Sprintf("Invalid filename(%s)", fileName)

		s.compriseMsg(w, body, http.StatusUnprocessableEntity)
		s.logMsg(s.warningLogger, body, http.StatusUnprocessableEntity)
		return

	} else {
		result, err := s.db.GetResValue(id)

		if s.reportUnreadyClient(w, id, result, err) {
			return

		} else if presignedURL := s.s3Client.GetViewFileURL(id, fileName); presignedURL != nil {
			body["Link"] = presignedURL.String()

			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			enc := json.NewEncoder(w)
			enc.SetEscapeHTML(false)
			enc.Encode(body)

			s.logMsg(s.debugLogger, body, http.StatusOK)

		} else {
			body["Error"] = fmt.Sprintf("Unable to find a file with such name %s", fileName)

			s.compriseMsg(w, body, http.StatusAccepted)
			s.logMsg(s.warningLogger, body, http.StatusAccepted)
		}
	}
}
