package server_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"net/http/httptest"
	"testing"

	config "web-service/src/config"
	http_server "web-service/src/server"
	utils "web-service/src/utils"
)

var server *http_server.Server

const (
	configPath = "../../config/main_config.json"
)

func UploadFilesActivityWorksInNormalUseCase(t *testing.T) {
	w := httptest.NewRecorder()
	content := []byte("This is a file with some content")
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)
	part, _ := writer.CreateFormFile("file_uploads", "sample_file.txt")
	part.Write(content)
	writer.Close()

	req := httptest.NewRequest("POST", "/", &buffer)
	req.Header.Set("Content-Type", "multipart/form-data; "+"boundary="+"\""+writer.Boundary()+"\"")
	server.UploadFilesHandler(w, req)
	result := w.Result()

	var unmarshaled map[string]interface{}
	body_bytes, _ := ioutil.ReadAll(result.Body)
	json.Unmarshal(body_bytes, &unmarshaled)
	t.Log(unmarshaled)

	if result.StatusCode != 200 {
		t.Error("Response status code should be 200")
	}
}

func UploadFilesActivityHandlesErrors(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	server.UploadFilesHandler(w, req)
	result := w.Result()
	if result.StatusCode == 200 {
		t.Error("Status code should not be 200")
	}
	if result.Body == nil {
		t.Error("Body should not be nil")
	}
	if value, ok := result.Header["Content-Type"]; ok {
		if len(value) != 1 {
			t.Error("Header slice should have length 1")
		}
		expected := "application/json; charset=utf-8"
		if value[0] != expected {
			t.Errorf("Content-Type should be %s", expected)
		}
	} else {
		t.Error("Header should contain Content-Type")
	}

	var body map[string]string
	json.NewDecoder(result.Body).Decode(&body)
	if _, ok := body["Error"]; !ok {
		t.Error("Header should contain Content-Type")
	}
}

func UploadFilesActivity(t *testing.T) {
	t.Run("UploadFiles", UploadFilesActivityHandlesErrors)
	t.Run("UploadFiles", UploadFilesActivityWorksInNormalUseCase)
}

func TestUploadFilesHandler(t *testing.T) {
	conf, err := config.NewConfig(configPath)
	if err != nil {
		t.Fatal("Unable to read program config from", configPath, err)
	}

	if err := utils.InitializeEmptyDir(conf.Internal.TempFilesDir); err != nil {
		t.Fatal("Unable to initialize temporary directory at", conf.Internal.TempFilesDir, err)
	}

	if err := utils.InitializeLogger(conf.Internal.LoggingDir); err != nil {
		t.Fatal("Unable to initialize logging at", conf.Internal.LoggingDir, err)
	}

	server = http_server.NewTestServer(conf)
	t.Run("UploadFiles", UploadFilesActivity)
}
