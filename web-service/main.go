package main

import (
	"log"

	config "web-service/src/config"
	s3support "web-service/src/s3support"
	http_server "web-service/src/server"
	stale_daemon "web-service/src/stale_daemon"
	containers "web-service/src/storage_container"
	nlp "web-service/src/text_similarity"
	utils "web-service/src/utils"
)

const (
	configPath = "config/main_config.json"
)

var (
	errorLogger *log.Logger
	debugLogger *log.Logger
)

func configure() {
	conf, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatalln("Unable to read program config from", configPath, err)
	}

	if err := utils.InitializeEmptyDir(conf.Internal.TempFilesDir); err != nil {
		log.Fatalln("Unable to initialize temporary directory at", conf.Internal.TempFilesDir, err)
	}

	if err := utils.InitializeLogger(conf.Internal.LoggingDir); err != nil {
		log.Fatalln("Unable to initialize logging at", conf.Internal.LoggingDir, err)
	}

	db, err := containers.NewDbClientService(conf.Internal.DbPath)
	if err != nil {
		errorLogger.Fatalln("Unable to open db at", conf.Internal.DbPath, err)
	}

	nlpCore := nlp.NewPyhonNLP(
		conf.Internal.PythonDifferenceScriptPath,
		conf.Internal.PythonSimilarityScriptPath,
	)

	s3Client, err := s3support.NewMinioService(&conf.Minio)
	if err != nil {
		errorLogger.Fatalln("Unable to initialize s3Client", err)
	}

	daemon := stale_daemon.NewDaemon(db, s3Client, conf.Internal.RefreshStaleDataPeriod)
	daemon.StartAsync()

	server := http_server.NewServer(conf, db, s3Client, nlpCore)
	server.Run()
}

func main() {
	configure()
}
