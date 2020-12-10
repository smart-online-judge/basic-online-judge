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
	MaxAllowedFilesSize        int64
	RefreshStaleDataPeriod     int64
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

func NewConfig(config_path string) (*Configuration, error) {
	content, err := ioutil.ReadFile(config_path)
	if err != nil {
		return nil, err
	}

	var conf Configuration
	if err = json.Unmarshal(content, &conf); err != nil {
		return nil, err
	}

	conf.Internal.PythonSimilarityScriptPath, err = utils.GetAbsolutePath(conf.Internal.PythonSimilarityScriptPath)
	if err != nil {
		return nil, err
	}

	conf.Internal.PythonDifferenceScriptPath, err = utils.GetAbsolutePath(conf.Internal.PythonDifferenceScriptPath)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}
