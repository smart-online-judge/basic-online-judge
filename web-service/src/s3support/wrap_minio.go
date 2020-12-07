package s3support

// This file provides the wrappers for use in Go code to save and load
// user files by their UUID v.4. We use min.io software, that is, the
// Amazon S3 compatible cloud object storage technology, to be able to
// horizontally scale this solution when we need more disk space.

// We belive that the metadata is already stored in SQLite database. So
// we simply work with the object storage, ignoring other concerns.

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"path"
	"strings"
	"time"

	config "web-service/src/config"
	utils "web-service/src/utils"

	guuid "github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	loggingPath = "logging/wrapMinIO.log"
)

var (
	minioClient *minio.Client
	bucketName  string

	warningLogger *log.Logger
	errorLogger   *log.Logger
	debugLogger   *log.Logger
	rootCtx       context.Context
)

func InitializeS3Support() {
	diagFile, err := utils.CreateFileIfNotExists(loggingPath)
	if err != nil {
		log.Fatalln("not able to initialize logger")
	}

	warningLogger = utils.GetLoggerPkgScoped("WARNING: ", diagFile)
	errorLogger = utils.GetLoggerPkgScoped("ERROR: ", diagFile)
	debugLogger = utils.GetLoggerPkgScoped("DEBUG: ", diagFile)
	rootCtx = context.Background()
	bucketName = config.Minio.BucketName

	minioClient, err = minio.New(config.Minio.ConnectionString, &minio.Options{
		Creds:  credentials.NewStaticV4(config.Minio.AccessKeyID, config.Minio.SecretAccessKey, ""),
		Secure: config.Minio.UseSSL},
	)
	if err != nil {
		errorLogger.Fatalln(err)
	}

	err = minioClient.MakeBucket(
		rootCtx,
		bucketName,

		// us-east-1 is used in local deployments
		minio.MakeBucketOptions{Region: "us-east-1"},
	)
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(rootCtx, bucketName)
		if errBucketExists == nil && exists {
			debugLogger.Printf("bucket %#v is already on the server\n", bucketName)
		} else {
			errorLogger.Fatalln(err)
		}
	} else {
		debugLogger.Println("created new bucket")
	}
}

func StoreFileByUUID(id guuid.UUID, file io.Reader, fileName string) error {
	objectName := path.Join(id.String(), "/", fileName)
	// - "application/octet-stream" means binary file, which we have
	// internally in Go here (io.Reader)
	// - "-1" means "unknown size"
	_, err := minioClient.PutObject(rootCtx, bucketName, objectName, file, -1, minio.PutObjectOptions{ContentType: "application/text", ContentEncoding: "utf-8"})

	if err != nil {
		errorLogger.Println(err)
	} else {
		debugLogger.Println("PutObject success")
	}
	return err
}

func UploadFsFileByUUID(id guuid.UUID, clientDir string, fileName string) error {
	objectName := path.Join(id.String(), "/", fileName)
	filePath := path.Join(clientDir, fileName)

	_, err := minioClient.FPutObject(
		rootCtx, bucketName, objectName, filePath,
		minio.PutObjectOptions{
			ContentType:     "application/text",
			ContentEncoding: "utf-8"},
	)

	if err != nil {
		errorLogger.Println(err)
	} else {
		debugLogger.Printf("PutObject %s to %s success\n", filePath, objectName)
	}
	return err
}

func DownloadFileByUUID(id guuid.UUID, fileName string) (io.Reader, error) {
	objectName := path.Join(id.String(), "/", fileName)
	v, err := minioClient.GetObject(rootCtx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		errorLogger.Println(err)
	}
	return v, err
}

func ListFilesByUUID(id guuid.UUID) []string {
	var res []string
	objectCh := minioClient.ListObjects(rootCtx, bucketName, minio.ListObjectsOptions{
		Prefix:    id.String(),
		Recursive: true,
	})
	for object := range objectCh {
		if object.Err != nil {
			errorLogger.Println(object.Err)
			return res
		}
		nextFileName := strings.Split(object.Key, "/")
		res = append(res, nextFileName[1])
	}
	return res
}

func GetViewFileURL(id guuid.UUID, fileName string) *url.URL {
	fileNames := ListFilesByUUID(id)
	ok := false
	for _, val := range fileNames {
		if fileName == val {
			ok = true
			break
		}
	}
	if ok == false {
		return nil
	}
	reqParams := make(url.Values)
	val := fmt.Sprintf("attachment; filename=\"%s\"", fileName)
	reqParams.Set("response-content-disposition", val)
	presignedURL, err := minioClient.PresignedGetObject(rootCtx, bucketName, path.Join(id.String(), fileName), time.Second*10*60, reqParams)
	if err != nil {
		errorLogger.Println(err)
		return nil
	}
	return presignedURL
}
