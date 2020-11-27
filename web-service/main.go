package main

import (
	"log"
	"net/http"
	handlers "web-service/src/api/controllers"
	containers "web-service/src/storage_container"
	utils "web-service/src/utils"
)

const (
	LoggingPath    = "logging/main_log.log"
	UploadFilesDir = "uploaded"
	DbPath         = "database/meta.db"
)

var (
	ErrorLogger *log.Logger
	DebugLogger *log.Logger
	container   *containers.DbClientContainer
)

func init() {
	utils.InitializeLogger(LoggingPath)

	ErrorLogger = utils.GetLogger("ERROR: ")
	DebugLogger = utils.GetLogger("DEBUG: ")

	container = containers.NewDB()

	container.Initialize(DbPath)
	DebugLogger.Println("Initialized database container in", DbPath)

	handlers.InitializeHandlersCommon()
	handlers.InitializeUploadFilesHandler(UploadFilesDir)
	DebugLogger.Println("Initialized upload Files handler in", UploadFilesDir)

	handlers.InitializeViewRoomHandler(container)
	DebugLogger.Println("Initialized view Room handler with", container)
}

func setupRoutes() {
	http.HandleFunc("/upload_files", handlers.UploadFilesHandler)
	http.HandleFunc("/view/", handlers.ViewRoomHandler)
}

func main() {
	setupRoutes()

	DebugLogger.Println("Starting fair online judge service on 8080 port")
	ErrorLogger.Fatal(http.ListenAndServe(":8080", nil))
}
