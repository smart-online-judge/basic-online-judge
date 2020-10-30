package daemon

import (
	"path"
	"path/filepath"
	"time"
	"log"
	"os"
	"encoding/json"
	"io/ioutil"
)

// Implement JSON Unmarshalling

// This struct uses struct tags, google it if you don't know
type StoredCalculations struct {
	StoredCalculations []StoredCalculation `json:"stored_calculations"`
}

type StoredCalculation struct {
	TimeStamp time.Time `json:"ts"`
	FileName string `json:"name"`
}

var (
	CommonLogger *log.Logger
	InfoLogger *log.Logger
)

func init() {
	fh, err := os.OpenFile("logs.txt", os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	CommonLogger = log.New(fh, "DAEMON_ERR: ", log.Ldate | log.Ltime | log.Lshortfile)
	InfoLogger = log.New(fh, "DAEMON_INFO", log.Ldate | log.Ltime | log.Lshortfile)
}

// This should be run once in every 5 minutes.
// Checks for data inside uploaded/ directory and if it does not
// find a valid timestamp in the `meta.json` file, removes a file
// TODO in init(), run a lambda function in a separate goroutine
// that waits for 5 minutes and calls this
func clean() {
	root := "uploaded"
	// if `root` does not exist, then exit
	if _, err := os.Stat(root); os.IsNotExist(err) {
		CommonLogger.Println(err)
		return
	}
	// Load meta.json into memory
	f, err := os.Open("meta.json")
	if err != nil {
		CommonLogger.Println(err)
	}
	defer f.Close()
	byteValue, _ := ioutil.ReadAll(f)
	var meta StoredCalculations
	json.Unmarshal(byteValue, &meta)
	// Clean the expired files
	err = filepath.Walk(root, func(_ string, info os.FileInfo, err error) error {
		InfoLogger.Printf("Checking if %s's storage has expired", info.Name())
		for _, storedCalculation := range meta.StoredCalculations {
			if storedCalculation.TimeStamp.Unix() < time.Now().Unix() {
				os.Remove(path.Join("uploaded", info.Name()))


			}
		}

		return nil
	})
}
