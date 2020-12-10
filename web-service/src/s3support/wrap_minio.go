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

type MinioService struct {
	minioClient *minio.Client
	bucketName  string

	warningLogger *log.Logger
	errorLogger   *log.Logger
	debugLogger   *log.Logger
	rootCtx       context.Context
}

func NewMinioService(conf *config.MinioConfiguration) (*MinioService, error) {
	minio_service := &MinioService{}
	err := minio_service.initialize(conf)
	return minio_service, err
}

func (self *MinioService) initialize(conf *config.MinioConfiguration) error {
	diagFile, err := utils.CreateFileIfNotExists(loggingPath)
	if err != nil {
		log.Fatalln("not able to initialize logger")
	}

	self.warningLogger = utils.GetLoggerPkgScoped("WARNING: ", diagFile)
	self.errorLogger = utils.GetLoggerPkgScoped("ERROR: ", diagFile)
	self.debugLogger = utils.GetLoggerPkgScoped("DEBUG: ", diagFile)
	self.rootCtx = context.Background()
	self.bucketName = conf.BucketName

	self.minioClient, err = minio.New(conf.ConnectionString, &minio.Options{
		Creds:  credentials.NewStaticV4(conf.AccessKeyID, conf.SecretAccessKey, ""),
		Secure: conf.UseSSL},
	)
	if err != nil {
		return err
	}

	err = self.minioClient.MakeBucket(
		self.rootCtx,
		self.bucketName,

		// us-east-1 is used in local deployments
		minio.MakeBucketOptions{Region: "us-east-1"},
	)
	if err != nil {
		exists, errBucketExists := self.minioClient.BucketExists(self.rootCtx, self.bucketName)
		if errBucketExists == nil && exists {
			self.debugLogger.Printf("bucket %#v is already on the server\n", self.bucketName)
		} else {
			return err
		}
	} else {
		self.debugLogger.Println("created new bucket")
	}
	return nil
}

func (self *MinioService) StoreFileByUUID(id guuid.UUID, file io.Reader, fileName string) error {
	objectName := path.Join(id.String(), "/", fileName)
	// - "application/octet-stream" means binary file, which we have
	// internally in Go here (io.Reader)
	// - "-1" means "unknown size"
	_, err := self.minioClient.PutObject(
		self.rootCtx, self.bucketName, objectName,
		file, -1, minio.PutObjectOptions{
			ContentType:     "application/text",
			ContentEncoding: "utf-8",
		},
	)

	if err != nil {
		self.errorLogger.Println(err)
	} else {
		self.debugLogger.Println("PutObject success")
	}
	return err
}

func (self *MinioService) UploadFsFileByUUID(id guuid.UUID, clientDir string, fileName string) error {
	objectName := path.Join(id.String(), "/", fileName)
	filePath := path.Join(clientDir, fileName)

	_, err := self.minioClient.FPutObject(
		self.rootCtx, self.bucketName, objectName, filePath,
		minio.PutObjectOptions{
			ContentType:     "application/text",
			ContentEncoding: "utf-8"},
	)

	if err != nil {
		self.errorLogger.Println(err)
	} else {
		self.debugLogger.Printf("PutObject %s to %s success\n", filePath, objectName)
	}
	return err
}

func (self *MinioService) DownloadFileByUUID(id guuid.UUID, fileName string) (io.Reader, error) {
	objectName := path.Join(id.String(), "/", fileName)
	v, err := self.minioClient.GetObject(
		self.rootCtx, self.bucketName,
		objectName, minio.GetObjectOptions{},
	)
	if err != nil {
		self.errorLogger.Println(err)
	}
	return v, err
}

func (self *MinioService) ListFilesByUUID(id guuid.UUID) []string {
	var res []string
	objectCh := self.minioClient.ListObjects(
		self.rootCtx, self.bucketName, minio.ListObjectsOptions{
			Prefix:    id.String(),
			Recursive: true},
	)
	for object := range objectCh {
		if object.Err != nil {
			self.errorLogger.Println(object.Err)
			return res
		}
		nextFileName := strings.Split(object.Key, "/")
		res = append(res, nextFileName[1])
	}
	return res
}

func (self *MinioService) GetViewFileURL(id guuid.UUID, fileName string) *url.URL {
	fileNames := self.ListFilesByUUID(id)
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
	presignedURL, err := self.minioClient.PresignedGetObject(
		self.rootCtx, self.bucketName, path.Join(id.String(), fileName), time.Second*10*60, reqParams)
	if err != nil {
		self.errorLogger.Println(err)
		return nil
	}
	return presignedURL
}

func (self *MinioService) RemoveFilesByPrefix(prefix string) {
	objectsCh := make(chan minio.ObjectInfo)

	// Send object names that are needed to be removed to objectsCh
	go func() {
		defer close(objectsCh)
		// List all objects from a bucket-name with a matching prefix.
		opts := minio.ListObjectsOptions{Prefix: prefix, Recursive: true}
		for object := range self.minioClient.ListObjects(self.rootCtx, self.bucketName, opts) {
			if object.Err != nil {
				self.errorLogger.Println(object.Err)
			}
			objectsCh <- object
		}
	}()

	// Call RemoveObjects API
	errorCh := self.minioClient.RemoveObjects(self.rootCtx, self.bucketName, objectsCh, minio.RemoveObjectsOptions{})

	// Print errors received from RemoveObjects API
	for e := range errorCh {
		self.errorLogger.Println("Failed to remove " + e.ObjectName + ", error: " + e.Err.Error())
	}
}
