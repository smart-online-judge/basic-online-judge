package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
)

func InitializeEmptyDir(dirPath string) error {
	dir, err := ioutil.ReadDir(dirPath)
	if os.IsNotExist(err) {
		if err = os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return err
		}
	} else {
		// Clear the directory content
		for _, d := range dir {
			os.RemoveAll(path.Join(dirPath, d.Name()))
		}
	}
	return nil
}

func CompriseMsg(w http.ResponseWriter, body map[string]interface{}, status int) error {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(body)
	return err
}

func LogMsg(logger *log.Logger, body map[string]interface{}, status int) {
	if status != http.StatusOK {
		logger.Println(body["Error"])
	} else {
		logger.Println(body["Result"])
	}
}
