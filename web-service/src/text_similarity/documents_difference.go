package text_similarity

import (
	"bytes"
	"os/exec"
	"strconv"

	config "web-service/src/config"
)

func GetFilesDifference(pipe_in bytes.Buffer, fileLen [2]int64, params ...string) (string, error) {
	var execute_path = config.Internal.PythonDifferenceScriptPath
	args := append([]string{execute_path}, strconv.FormatInt(fileLen[0], 10))
	args = append(args, strconv.FormatInt(fileLen[1], 10))
	args = append(args, params...)
	var pipe_out, pipe_err bytes.Buffer

	// rely on environment variable for `python`
	cmd := exec.Command("python", args...)
	cmd.Stdin = &pipe_in
	cmd.Stdout = &pipe_out
	cmd.Stderr = &pipe_err

	if err := cmd.Run(); err != nil {
		return "", NewPythonError(pipe_err.String())
	}

	return pipe_out.String(), nil
}
