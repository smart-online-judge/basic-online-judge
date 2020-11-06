package text_similarity

import (
	"bytes"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	python_script_path      = "api/python_wrappers/similarity_module.py"
	python_interpreter_path = "/usr/bin/python3"
)

func GetPairwiseSimilarity(input_path string, args ...string) ([][]float32, error) {
	var err error
	var cur_path string
	if cur_path, err = filepath.Abs(filepath.Dir(os.Args[0])); err != nil {
		return nil, err
	}

	var execute_path = python_script_path
	if !strings.HasPrefix(python_script_path, string(os.PathSeparator)) {
		execute_path = path.Join(cur_path, execute_path)
	}

	if !strings.HasPrefix(input_path, string(os.PathSeparator)) {
		input_path = path.Join(cur_path, input_path)
	}

	args = append([]string{execute_path}, args...)
	var pipe_out bytes.Buffer
	cmd := exec.Command(python_interpreter_path, append(args, input_path)...)
	cmd.Stdout = &pipe_out
	cmd.Stderr = os.Stderr

	if err = cmd.Run(); err != nil {
		return nil, err
	}

	out_arr := strings.Split(pipe_out.String(), ",")
	mat_size, _ := strconv.Atoi(out_arr[0])

	res_mat := make([][]float32, mat_size)
	for i := range res_mat {
		res_mat[i] = make([]float32, mat_size)
	}

	for idx, upper_b := 1, mat_size*mat_size; idx < upper_b; idx++ {
		val, _ := strconv.ParseFloat(out_arr[idx], 32)
		res_mat[idx/mat_size][idx%mat_size] = float32(val)
	}

	return res_mat, nil
}
