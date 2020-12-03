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

	if req.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		warningLogger.Printf("%s request type isn't supported\n", req.Method)
		return
	}

	var err error
	if err = req.ParseMultipartForm(int64(maxAllowedFilesSize)); nil != err {
		http.Error(w, "Unable to parse input files", http.StatusPreconditionFailed)
		errorLogger.Println(err)
		return
	}

	id := guuid.New()
	err = os.Mkdir(path.Join(uploadFilesDir, id.String()), os.ModePerm)
	if err != nil {
		http.Error(w, "Unable to store input files", http.StatusInternalServerError)
		errorLogger.Println(err)
		return
	}

	fhs := req.MultipartForm.File["file_uploads"]
	for _, fh := range fhs {
		var infile multipart.File
		if infile, err = fh.Open(); nil != err {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			errorLogger.Println(err)
			return
		}
		defer infile.Close()

		var outfile *os.File
		if outfile, err = os.Create(path.Join(uploadFilesDir, id.String(), fh.Filename)); nil != err {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			errorLogger.Println(err)
			return
		}

		if _, err = io.Copy(outfile, infile); nil != err {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			errorLogger.Println(err)
			return
		}
	}
	go prepareViewForUUID(id)
	// Report a link to the personal room
	fmt.Fprintf(w, "You can view the result of the file analyses at the view/%v\n", id)
}
