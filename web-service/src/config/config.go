package config

import (
	"encoding/json"
	"io/ioutil"

	utils "web-service/src/utils"
)

type ServerConfiguration struct {
	Port string
}

type InternalConfiguration struct {
	TempFilesDir               string
	UploadFilesDir             string
	MaxAllowedFilesSize        int
	LoggingDir                 string
	DbPath                     string
	PythonSimilarityScriptPath string
	PythonDifferenceScriptPath string
}

type MinioConfiguration struct {
	ConnectionString string
	AccessKeyID      string
	SecretAccessKey  string
	UseSSL           bool
	BucketName       string
}

type Configuration struct {
	Server   ServerConfiguration
	Internal InternalConfiguration
	Minio    MinioConfiguration
}

var (
	Server   ServerConfiguration
	Internal InternalConfiguration
	Minio    MinioConfiguration
)

func ReadConfig(config_path string) error {
	content, err := ioutil.ReadFile(config_path)
	if err != nil {
		return err
	}

	var conf Configuration
	if err = json.Unmarshal(content, &conf); err != nil {
		return err
	}

	Server = conf.Server
	Internal = conf.Internal
	Minio = conf.Minio

	Internal.PythonSimilarityScriptPath, err = utils.GetAbsolutePath(Internal.PythonSimilarityScriptPath)
	if err != nil {
		return err
	}

	Internal.PythonDifferenceScriptPath, err = utils.GetAbsolutePath(Internal.PythonDifferenceScriptPath)
	if err != nil {
		return err
	}

	return nil
}
