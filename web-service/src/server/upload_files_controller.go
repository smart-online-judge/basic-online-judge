package server

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"

	guuid "github.com/google/uuid"
)

func (s *Server) UploadFilesHandler(w http.ResponseWriter, req *http.Request) {
	s.debugLogger.Println("uploadFiles Endpoint hit")
	body := make(map[string]interface{})

	if req.Method != "POST" {
		body["Error"] = fmt.Sprintf("%s request type isn't supported", req.Method)
		s.compriseMsg(w, body, http.StatusMethodNotAllowed)
		s.logMsg(s.warningLogger, body, http.StatusMethodNotAllowed)
		return
	}

	var err error
	if err = req.ParseMultipartForm(s.conf.Internal.MaxAllowedFilesSize); err != nil {
		body["Message"] = "Unable to parse input files"
		body["Error"] = err.Error()

		s.compriseMsg(w, body, http.StatusPreconditionFailed)
		s.logMsg(s.errorLogger, body, http.StatusPreconditionFailed)
		return
	}

	id := guuid.New()
	if err = os.Mkdir(path.Join(s.conf.Internal.UploadFilesDir, id.String()), os.ModePerm); err != nil {
		body["Message"] = "Unable to store input files"
		body["Error"] = err.Error()

		s.compriseMsg(w, body, http.StatusInternalServerError)
		s.logMsg(s.errorLogger, body, http.StatusInternalServerError)
		return
	}

	fhs := req.MultipartForm.File["file_uploads"]
	for _, fh := range fhs {
		var infile multipart.File
		if infile, err = fh.Open(); err != nil {
			body["Message"] = http.StatusText(http.StatusInternalServerError)
			body["Error"] = err.Error()

			s.compriseMsg(w, body, http.StatusInternalServerError)
			s.logMsg(s.errorLogger, body, http.StatusInternalServerError)
			return
		}
		defer infile.Close()

		var outfile *os.File
		if outfile, err = os.Create(path.Join(s.conf.Internal.UploadFilesDir, id.String(), fh.Filename)); err != nil {
			body["Message"] = http.StatusText(http.StatusInternalServerError)
			body["Error"] = err.Error()

			s.compriseMsg(w, body, http.StatusInternalServerError)
			s.logMsg(s.errorLogger, body, http.StatusInternalServerError)
			return
		}

		if _, err = io.Copy(outfile, infile); err != nil {
			body["Message"] = http.StatusText(http.StatusInternalServerError)
			body["Error"] = err.Error()

			s.compriseMsg(w, body, http.StatusInternalServerError)
			s.logMsg(s.errorLogger, body, http.StatusInternalServerError)
			return
		}
	}
	go s.prepareViewForUUID(id)
	// Report a link to the personal room
	body["Message"] = fmt.Sprintf("You can view the result of the file analyses at the view %s", id)
	body["Result"] = id

	s.compriseMsg(w, body, http.StatusOK)
	s.logMsg(s.debugLogger, body, http.StatusOK)
}
