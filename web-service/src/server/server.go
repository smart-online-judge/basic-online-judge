package server

import (
	"log"
	"net/http"
	"os"

	config "web-service/src/config"
	s3support "web-service/src/s3support"
	containers "web-service/src/storage_container"
	nlp "web-service/src/text_similarity"
	utils "web-service/src/utils"
)

type Server struct {
	conf *config.Configuration

	db       containers.ClientStorageInterface
	s3Client *s3support.MinioService
	nlpCore  nlp.NlpModuleInterface

	warningLogger *log.Logger
	errorLogger   *log.Logger
	debugLogger   *log.Logger
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/upload_files", s.UploadFilesHandler)
	mux.HandleFunc("/api/view/", s.ViewRoomHandler)
	mux.HandleFunc("/api/link", s.GetFileLinkById)
	mux.HandleFunc("/api/cmp_files", s.CompareFilesHandler)

	return mux
}

func (s *Server) Run() {
	httpServer := &http.Server{
		Addr:    ":" + s.conf.Server.Port,
		Handler: s.Handler(),
	}

	s.errorLogger.Println(httpServer.ListenAndServe())
}

func NewServer(
	conf *config.Configuration,
	container containers.ClientStorageInterface,
	s3Client *s3support.MinioService,
	nlpCore nlp.NlpModuleInterface) *Server {

	s := &Server{
		conf:          conf,
		db:            container,
		s3Client:      s3Client,
		warningLogger: utils.WarningLogger,
		debugLogger:   utils.DebugLogger,
		errorLogger:   utils.ErrorLogger,
		nlpCore:       nlpCore,
	}

	if err := os.Mkdir(conf.Internal.UploadFilesDir, os.ModePerm); err != nil {
		s.errorLogger.Fatalln(err)
	}

	return s
}

func NewTestServer(
	conf *config.Configuration) *Server {

	s := &Server{
		conf:          conf,
		db:            nil,
		s3Client:      nil,
		warningLogger: utils.WarningLogger,
		debugLogger:   utils.DebugLogger,
		errorLogger:   utils.ErrorLogger,
		nlpCore:       nil,
	}

	if err := os.Mkdir(conf.Internal.UploadFilesDir, os.ModePerm); err != nil {
		s.errorLogger.Fatalln(err)
	}

	return s
}
