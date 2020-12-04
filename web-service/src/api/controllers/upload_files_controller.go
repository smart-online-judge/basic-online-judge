package api

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"

	guuid "github.com/google/uuid"
	config "web-service/src/config"
)

var (
	maxAllowedFilesSize int
)

func initializeUploadFilesController() {
	maxAllowedFilesSize = config.Internal.MaxAllowedFilesSize

	var err error
	if _, err = os.Stat(uploadFilesDir); os.IsNotExist(err) {
		err = os.MkdirAll(uploadFilesDir, os.ModePerm)
	}
	if err != nil {
		errorLogger.Fatalln(err)
	}
}

func UploadFilesHandler(w http.ResponseWriter, req *http.Request) {
	debugLogger.Println("uploadFiles Endpoint hit")
	body := make(map[string]interface{})

	if req.Method != "POST" {
		body["Error"] = fmt.Sprintf("%s request type isn't supported", req.Method)
		compriseMsg(w, body, http.StatusMethodNotAllowed)
		logMsg(warningLogger, body, http.StatusMethodNotAllowed)
		return
	}

	var err error
	if err = req.ParseMultipartForm(int64(maxAllowedFilesSize)); err != nil {
		body["Message"] = "Unable to parse input files"
		body["Error"] = err.Error()

		compriseMsg(w, body, http.StatusPreconditionFailed)
		logMsg(errorLogger, body, http.StatusPreconditionFailed)
		return
	}

	id := guuid.New()
	if err = os.Mkdir(path.Join(uploadFilesDir, id.String()), os.ModePerm); err != nil {
		body["Message"] = "Unable to store input files"
		body["Error"] = err.Error()

		compriseMsg(w, body, http.StatusInternalServerError)
		logMsg(errorLogger, body, http.StatusInternalServerError)
		return
	}

	fhs := req.MultipartForm.File["file_uploads"]
	for _, fh := range fhs {
		var infile multipart.File
		if infile, err = fh.Open(); err != nil {
			body["Message"] = http.StatusText(http.StatusInternalServerError)
			body["Error"] = err.Error()

			compriseMsg(w, body, http.StatusInternalServerError)
			logMsg(errorLogger, body, http.StatusInternalServerError)
			return
		}
		defer infile.Close()

		var outfile *os.File
		if outfile, err = os.Create(path.Join(uploadFilesDir, id.String(), fh.Filename)); err != nil {
			body["Message"] = http.StatusText(http.StatusInternalServerError)
			body["Error"] = err.Error()

			compriseMsg(w, body, http.StatusInternalServerError)
			logMsg(errorLogger, body, http.StatusInternalServerError)
			return
		}

		if _, err = io.Copy(outfile, infile); err != nil {
			body["Message"] = http.StatusText(http.StatusInternalServerError)
			body["Error"] = err.Error()

			compriseMsg(w, body, http.StatusInternalServerError)
			logMsg(errorLogger, body, http.StatusInternalServerError)
			return
		}
	}
	go prepareViewForUUID(id)
	// Report a link to the personal room
	body["Message"] = fmt.Sprintf("You can view the result of the file analyses at the view %s", id)
	body["Result"] = id

	compriseMsg(w, body, http.StatusOK)
	logMsg(debugLogger, body, http.StatusOK)
}
