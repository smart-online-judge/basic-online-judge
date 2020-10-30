package main

import (
	"fmt"
	"io"
	//"strconv"
	"log"
	"mime/multipart"
	"net/http"
	//"net/url"
	"os"
	"path"
	// TODO Make sure "NLP" is seen here.
	// I want to see where the module comes from, not guess.
	// See how guuid is imported
	"project/python_wrappers"
	guuid "github.com/google/uuid"
)

var (
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
	UUIDInMemCache map[guuid.UUID]bool
	// TODO remove this, leave UUIDResult
	UUIDIsReadyForView map[guuid.UUID]bool
	UUIDResult map[guuid.UUID][][]float32
)

func init() {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat("uploaded"); os.IsNotExist(err) {
		// TODO 0777 is too much privilege, but it does not work with 0666.
		err := os.Mkdir("uploaded", 0777)
		if err != nil {
			log.Fatal(err)
		}
	}

	WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	UUIDInMemCache = make(map[guuid.UUID]bool)
	UUIDIsReadyForView = make(map[guuid.UUID]bool)
	UUIDResult = make(map[guuid.UUID][][]float32)
}

func prepareViewForUUID(id guuid.UUID) {
	// TODO think about concurrency
	UUIDInMemCache[id] = true
	res, err := nlp.ComputePairwiseSimilarity("uploaded", "--external")
	if err != nil {
		ErrorLogger.Println(err)
	}
	UUIDIsReadyForView[id] = true
	UUIDResult[id] = res
}

// Handlers

func viewSimilarity(w http.ResponseWriter, req *http.Request) {
	fmt.Println("viewSimilarity End Point hit")
	// Retrieve view id
	var id_ns, ok = req.URL.Query()["id"]
	if !ok {
		ErrorLogger.Println("Invalid url format")
		return
	}
	id_str := id_ns[0]
	id, err := guuid.Parse(id_str)
	if err != nil {
		ErrorLogger.Println("Invalid UUID value")
		return
	}
	// If we are not ready, deny of service.
	if _, err := UUIDIsReadyForView[id]; err == true {
		fmt.Printf("Result for %s is not yet ready", id);
		return
	}
	// Serve the result.
	fmt.Println("%b", UUIDResult[id])
}

func uploadFiles(w http.ResponseWriter, req *http.Request) {
	fmt.Println("File Upload Endpoint Hit")
	var err error

	if req.Method == "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		WarningLogger.Println("Get request not supported")
		return
	}

	const _32K = (1 << 10) * 32
	if err = req.ParseMultipartForm(_32K); nil != err {
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
		if outfile, err = os.Create(path.Join("uploaded", fh.Filename)); nil != err {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			ErrorLogger.Println(err)
			return
		}

		// 32K buffer copy
		if _, err = io.Copy(outfile, infile); nil != err {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			ErrorLogger.Println(err)
			return
		}

		// fmt.Printf("Uploaded File: %+v\n", fh.Filename)
		// fmt.Printf("File Size: %+v\n", fh.Size)
		// fmt.Printf("MIME Header: %+v\n", fh.Header)
	}

	/// Report a link to the personal room

	// TODO Fix concurrent access/modification
	id := guuid.New()
	for {
		if _, ok := UUIDInMemCache[id]; ok {
			id = guuid.New()
			continue
		}
		break
	}
	go prepareViewForUUID(id)

	fmt.Fprint(w, "%s", id.String())
}

// Necessary

func setupRoutes() {
	http.HandleFunc("/upload", uploadFiles)
	http.HandleFunc("/view", viewSimilarity)
	http.ListenAndServe(":8080", nil)
}

func main() {
	fmt.Println("Hello World")
	setupRoutes()
}
