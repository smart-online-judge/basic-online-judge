package main

import (
	"log"
	"net/http"
	handlers "web-service/api/handlers"
	utils "web-service/api/utils"
)

const (
	LoggingPath    = "api/logging/main_log.log"
	UploadFilesDir = "uploaded"
)

var (
	ErrorLogger *log.Logger
	DebugLogger *log.Logger
)

func init() {
	utils.InitializeLogger(LoggingPath)

	ErrorLogger = utils.GetLogger("ERROR: ")
	DebugLogger = utils.GetLogger("DEBUG: ")
}

func setupRoutes() {
	handlers.PrepareHandlersCommon()
	handlers.PrepareUploadFilesHandler(UploadFilesDir)
	handlers.PrepareViewRoomHandler()

	http.HandleFunc("/upload_files", handlers.UploadFilesHandler)
	http.HandleFunc("/view", handlers.ViewRoomHandler)
}

func main() {
	setupRoutes()

	DebugLogger.Println("Starting fair online judge service on 8080 port")
	ErrorLogger.Fatal(http.ListenAndServe(":8080", nil))
}
