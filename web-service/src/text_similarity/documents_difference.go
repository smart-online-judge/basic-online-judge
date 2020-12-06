package text_similarity

import (
	"bytes"
	//"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	config "web-service/src/config"
)

func GetFilesDifference(pipe_in bytes.Buffer, fileLen [2]int64, params ...string) (string, error) {
	var err error
	var cur_path string
	if cur_path, err = filepath.Abs(filepath.Dir(os.Args[0])); err != nil {
		return "", err
	}

	var execute_path = config.Internal.PythonDifferenceScriptPath
	if !strings.HasPrefix(execute_path, string(os.PathSeparator)) {
		execute_path = path.Join(cur_path, execute_path)
	}

	args := append([]string{execute_path}, strconv.FormatInt(fileLen[0], 10))
	args = append(args, strconv.FormatInt(fileLen[1], 10))
	args = append(args, params...)
	var pipe_out, pipe_err bytes.Buffer

	// rely on environment variable for `python`
	cmd := exec.Command("python", args...)
	cmd.Stdin = &pipe_in
	cmd.Stdout = &pipe_out
	cmd.Stderr = &pipe_err

	if err = cmd.Run(); err != nil {
		return "", NewPythonError(pipe_err.String())
	}

	return pipe_out.String(), nil
}
