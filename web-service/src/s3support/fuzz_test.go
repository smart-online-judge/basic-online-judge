package s3support_test

import (
	"testing"
	"bytes"

	fuzz "github.com/google/gofuzz"
	guuid "github.com/google/uuid"

	config "web-service/src/config"
	s3support "web-service/src/s3support"
	utils "web-service/src/utils"

)

var minio *s3support.MinioService

const (
	configPath = "../../config/main_config.json"
)


func TestSetup(t *testing.T) {
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

	minio, err = s3support.NewMinioService(&conf.Minio)
	if err != nil {
		t.Fatal("Unable to initialize s3Client", err)
	}

	t.Run("Fuzz", TestFuzz)
}

// This method is an example of how to use fuzzing now, it might be
// useful later.
func TestStoreFileByUUIDDoesNotPanicWithRandomData(t *testing.T) {
	// Configre the Fuzzer, if needed
	f := fuzz.New()

	var raw_bytes []byte
	f.Fuzz(&raw_bytes)

	buffer := bytes.NewBuffer(raw_bytes)

	var id guuid.UUID
	f.Fuzz(&id)

	var fileName string
	f.Fuzz(&fileName)

	err := minio.StoreFileByUUID(id, buffer, fileName)
	if err != nil {
		panic(err)
	}
}
