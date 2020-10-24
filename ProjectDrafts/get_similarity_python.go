package main

import (
    "log"
    "os"
	"path"
	"os/exec"
	"fmt"
	"strings"
	"bytes"
	"path/filepath"
	"strconv"
)

func handle_error(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	if len(os.Args) == 1 {
        return
	}

	cur_path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	handle_error(err)

	if !strings.HasPrefix(os.Args[1], string(os.PathSeparator)) {
		os.Args[1] = path.Join(cur_path, os.Args[1])
	}
	var pipe_out bytes.Buffer

    cmd := exec.Command("/usr/bin/python3", os.Args[1:]...)
    cmd.Stdout = &pipe_out
	cmd.Stderr = os.Stderr
	
	err = cmd.Run()
	handle_error(err)

	out_arr := strings.Split(pipe_out.String(), ",")
	mat_size, _ := strconv.Atoi(out_arr[0])

	res_mat := make([][]float32, mat_size)
	for i := range res_mat {
		res_mat[i] = make([]float32, mat_size)
	}

	for index, str := range out_arr[1:] {
		val, _ := strconv.ParseFloat(str, 32)
		res_mat[index / mat_size][index % mat_size] = float32(val)
	}
	fmt.Println(res_mat)
}