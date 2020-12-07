package text_similarity

import (
	"bytes"
	"fmt"
	"os/exec"

	config "web-service/src/config"
	utils "web-service/src/utils"
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

func GetPairwiseSimilarity(input_path string, args ...string) (string, error) {
	var err error
	if input_path, err = utils.GetAbsolutePath(input_path); err != nil {
		return "", err
	}

	var execute_path = config.Internal.PythonSimilarityScriptPath
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
