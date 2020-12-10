package text_similarity

import (
	"bytes"
	"os/exec"
	"strconv"

	utils "web-service/src/utils"
)

type PythonNLP struct {
	diffScriptPath       string
	similarityScriptPath string
}

func NewPyhonNLP(diff_path, similarity_path string) *PythonNLP {
	return &PythonNLP{diff_path, similarity_path}
}

func (self *PythonNLP) GetPairwiseSimilarity(input_path string, args ...string) (string, error) {
	var err error
	if input_path, err = utils.GetAbsolutePath(input_path); err != nil {
		return "", err
	}

	args = append([]string{self.similarityScriptPath}, args...)

	var pipe_out, pipe_err bytes.Buffer

	cmd := exec.Command("python", append(args, input_path)...)
	cmd.Stdout = &pipe_out
	cmd.Stderr = &pipe_err

	if err = cmd.Run(); err != nil {
		return "", NewPythonError(pipe_err.String())
	}

	return pipe_out.String(), nil
}

func (self *PythonNLP) GetFilesDifference(pipe_in bytes.Buffer, fileLen [2]int64, params ...string) (string, error) {
	args := append([]string{self.diffScriptPath}, strconv.FormatInt(fileLen[0], 10))
	args = append(args, strconv.FormatInt(fileLen[1], 10))
	args = append(args, params...)
	var pipe_out, pipe_err bytes.Buffer

	cmd := exec.Command("python", args...)
	cmd.Stdin = &pipe_in
	cmd.Stdout = &pipe_out
	cmd.Stderr = &pipe_err

	if err := cmd.Run(); err != nil {
		return "", NewPythonError(pipe_err.String())
	}

	return pipe_out.String(), nil
}
