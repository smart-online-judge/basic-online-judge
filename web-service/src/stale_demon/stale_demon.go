package stale_demon

import (
	"log"
	"time"

	s3support "web-service/src/s3support"
	containers "web-service/src/storage_container"
	utils "web-service/src/utils"
)

var (
	db          containers.ClientContainer
	runPeriod   int64
	errorLogger *log.Logger
	debugLogger *log.Logger
)

func InitializeDemon(container containers.ClientContainer, period int64) {
	db = container
	runPeriod = period

	errorLogger = utils.ErrorLogger
	debugLogger = utils.DebugLogger
}

func removeStaleData() {
	debugLogger.Println("Started removeStaleData() function")
	ids, err := db.GetRemoveStaleClients(runPeriod)

	if err != nil {
		errorLogger.Println("Error encountered when removing stale data: ", err)
		return
	}

	for _, id := range ids {
		s3support.RemoveFilesByPrefix(id)
	}
	debugLogger.Println("Demon removed", len(ids), "clients")
}

func StartAsync() *time.Ticker {
	debugLogger.Printf("Started stale demon with %ds period\n", runPeriod)

	ticker := time.NewTicker(time.Duration(runPeriod) * time.Second)
	go func() {
		for range ticker.C {
			removeStaleData()
		}
	}()
	return ticker
}
