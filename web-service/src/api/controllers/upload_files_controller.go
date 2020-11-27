package api

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"

	config "web-service/src/config"

	guuid "github.com/google/uuid"
)

func UploadFilesHandler(w http.ResponseWriter, req *http.Request) {
	DebugLogger.Println("uploadFiles Endpoint hit")

	if req.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		WarningLogger.Printf("%s request type isn't supported\n", req.Method)
		return
	}

	var err error
	const _20M = (1 << 20) * 20
	if err = req.ParseMultipartForm(_20M); nil != err {
		http.Error(w, "507 - Maximum upload size limit exceeded!", http.StatusInsufficientStorage)
		return
	}

	id := guuid.New()
	os.Mkdir(path.Join(config.UploadFilesDir, id.String()), os.ModePerm)

	fhs := req.MultipartForm.File["givenFiles"]
	for _, fh := range fhs {
		var infile multipart.File
		if infile, err = fh.Open(); nil != err {
			ErrorLogger.Println(err)
			return
		}
		defer infile.Close()

		var outfile *os.File
		if outfile, err = os.Create(path.Join(config.UploadFilesDir, id.String(), fh.Filename)); nil != err {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			ErrorLogger.Println(err)
			return
		}

		// 20M buffer copy
		if _, err = io.Copy(outfile, infile); nil != err {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			ErrorLogger.Println(err)
			return
		}
	}
	go prepareViewForUUID(id)
	// Report a link to the personal room
	fmt.Fprintf(w, "View room id= %v", id)
}
