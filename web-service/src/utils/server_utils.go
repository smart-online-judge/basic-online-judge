package utils

import (
	"io/ioutil"
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
