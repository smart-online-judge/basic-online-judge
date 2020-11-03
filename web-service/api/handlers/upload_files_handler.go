package handlers

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	guuid "github.com/google/uuid"
)

var (
	UploadFilesDir  string
)

func PrepareUploadFilesHandler(upload_files_dir string) {

	if _, err := os.Stat(upload_files_dir); os.IsNotExist(err) {
		if err := os.MkdirAll(upload_files_dir, 0777); err != nil {
			ErrorLogger.Fatal(err)
		}
	}

	UploadFilesDir = upload_files_dir
}

func UploadFilesHandler(w http.ResponseWriter, req *http.Request) {
	DebugLogger.Println("uploadFiles Endpoint hit")

	if req.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		WarningLogger.Printf("%s request type isn't supported\n", req.Method)
		return
	}

	var err error
	const _5M = (1 << 20) * 5
	if err = req.ParseMultipartForm(_5M); nil != err {
		http.Error(w, "507 - Maximum upload size limit exceeded!", http.StatusInsufficientStorage)
		return
	}

	fhs := req.MultipartForm.File["givenFiles"]
	for _, fh := range fhs {
		var infile multipart.File
		if infile, err = fh.Open(); nil != err {
			ErrorLogger.Println(err)
			return
		}
		defer infile.Close()

		var outfile *os.File
		if outfile, err = os.Create(path.Join(UploadFilesDir, fh.Filename)); nil != err {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			ErrorLogger.Println(err)
			return
		}

		// 5M buffer copy
		if _, err = io.Copy(outfile, infile); nil != err {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			ErrorLogger.Println(err)
			return
		}
	}

	id := guuid.New()
	go prepareViewForUUID(id)
	// Report a link to the personal room
	fmt.Fprintf(w, "View room id= %v", id)
}