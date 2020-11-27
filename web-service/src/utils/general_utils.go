package utils

import "os"

type Pair struct {
	First, Second interface{}
}

// Create a File if it does not exist
func CreateFileIfNotExists(path string) (*os.File, error) {
	fd, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	return fd, nil
}
