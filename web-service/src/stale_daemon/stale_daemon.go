package stale_daemon

import (
	"log"
	"time"

	s3support "web-service/src/s3support"
	containers "web-service/src/storage_container"
	utils "web-service/src/utils"
)

type StaleDaemon struct {
	db          containers.ClientStorageInterface
	s3Client    *s3support.MinioService
	runPeriod   int64
	ticker      *time.Ticker
	errorLogger *log.Logger
	debugLogger *log.Logger
}

func NewDaemon(
	container containers.ClientStorageInterface,
	s3Client *s3support.MinioService,
	period int64) *StaleDaemon {

	daemon := &StaleDaemon{
		db:          container,
		s3Client:    s3Client,
		runPeriod:   period,
		errorLogger: utils.ErrorLogger,
		debugLogger: utils.DebugLogger,
	}
	return daemon
}

func (self *StaleDaemon) removeStaleData() {
	self.debugLogger.Println("Started removeStaleData() function")
	ids, err := self.db.GetRemoveStaleClients(self.runPeriod)

	if err != nil {
		self.errorLogger.Println("Error encountered when removing stale data: ", err)
		return
	}

	for _, id := range ids {
		self.s3Client.RemoveFilesByPrefix(id)
	}
	self.debugLogger.Println("Daemon removed", len(ids), "clients")
}

func (self *StaleDaemon) StartAsync() {
	self.debugLogger.Printf("Started stale daemon with %ds period\n", self.runPeriod)

	self.ticker = time.NewTicker(time.Duration(self.runPeriod) * time.Second)
	go func() {
		for ; true; <-self.ticker.C {
			self.removeStaleData()
		}
	}()
}
