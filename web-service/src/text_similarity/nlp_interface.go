package text_similarity

import (
	"bytes"
	"fmt"
)

type PythonInternalError struct {
	ErrMsg string
}

func NewPythonError(err_msg string) error {
	return &PythonInternalError{err_msg}
}

func (e *PythonInternalError) Error() string {
	return fmt.Sprintf("NLP error: %v", e.ErrMsg)
}

type NlpModuleInterface interface {
	GetPairwiseSimilarity(input_path string, args ...string) (string, error)
	GetFilesDifference(pipe_in bytes.Buffer, fileLen [2]int64, params ...string) (string, error)
}
