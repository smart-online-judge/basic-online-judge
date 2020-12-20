package utils

import (
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Pair struct {
	First, Second interface{}
}

// Create a File if it does not exist
func CreateFileIfNotExists(pathStr string) (*os.File, error) {
	dir, _ := path.Split(pathStr)
	if _, err := os.Stat(pathStr); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}
	fd, err := os.OpenFile(pathStr, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}
	return fd, nil
}

func GetAbsolutePath(input_path string) (string, error) {
	var err error
	var cur_path string
	if cur_path, err = filepath.Abs(filepath.Dir(os.Args[0])); err != nil {
		return "", err
	}

	if !strings.HasPrefix(input_path, string(os.PathSeparator)) {
		input_path = path.Join(cur_path, input_path)
	}
	return input_path, nil
}
