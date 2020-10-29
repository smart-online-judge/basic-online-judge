package main

import (
	"fmt"
	"io"
	//"strconv"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"project/python_wrappers"
)

var (
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
)

func init() {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func uploadFile(w http.ResponseWriter, req *http.Request) {
	fmt.Println("File Upload Endpoint Hit")
	var err error

	if req.Method == "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		WarningLogger.Println("Get request not supported")
		return
	}

	const _32M = (1 << 20) * 32
	if err = req.ParseMultipartForm(_32M); nil != err {
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

		// 32M buffer copy
		if _, err = io.Copy(outfile, infile); nil != err {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			ErrorLogger.Println(err)
			return
		}

		fmt.Printf("Uploaded File: %+v\n", fh.Filename)
		fmt.Printf("File Size: %+v\n", fh.Size)
		fmt.Printf("MIME Header: %+v\n", fh.Header)
		fmt.Fprintf(w, "Successfully Uploaded File\n")
	}

	if res, err := nlp.ComputePairwiseSimilarity("uploaded", "--external"); err != nil {
		ErrorLogger.Println(err)
	} else {
		fmt.Fprintf(w, "%v\n", res)
	}
}

func setupRoutes() {
	http.HandleFunc("/upload", uploadFile)
	http.ListenAndServe(":8080", nil)
}

func main() {
	fmt.Println("Hello World")
	setupRoutes()

	//res, err := nlp.ComputePairwiseSimilarity("uploaded", "--external")
	//fmt.Println(res, err)
}
