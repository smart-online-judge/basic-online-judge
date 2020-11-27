package text_similarity

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

const (
	python_script_path = "src/python/get_documents_similarity.py"
)

type PythonInternalError struct {
	ErrMsg string
}

func NewPythonError(err_msg string) error {
	return &PythonInternalError{err_msg}
}

func (e *PythonInternalError) Error() string {
	return fmt.Sprintf("python scipt error: %v", e.ErrMsg)
}

func GetPairwiseSimilarity(input_path string, args ...string) (string, error) {
	var err error
	var cur_path string
	if cur_path, err = filepath.Abs(filepath.Dir(os.Args[0])); err != nil {
		return "", err
	}

	var execute_path = python_script_path
	if !strings.HasPrefix(python_script_path, string(os.PathSeparator)) {
		execute_path = path.Join(cur_path, execute_path)
	}

	if !strings.HasPrefix(input_path, string(os.PathSeparator)) {
		input_path = path.Join(cur_path, input_path)
	}

	args = append([]string{execute_path}, args...)
	var pipe_out, pipe_err bytes.Buffer
	// rely on environment variable for `python`
	cmd := exec.Command("python", append(args, input_path)...)
	cmd.Stdout = &pipe_out
	cmd.Stderr = &pipe_err

	if err = cmd.Run(); err != nil {
		return "", NewPythonError(pipe_err.String())
	}

	return pipe_out.String(), nil
}
