package s3support

// This file provides the wrappers for use in Go code to save and load
// user files by their UUID v.4. We use min.io software, that is, the
// Amazon S3 compatible cloud object storage technology, to be able to
// horizontally scale this solution when we need more disk space.

// We belive that the metadata is already stored in SQLite database. So
// we simply work with the object storage, ignoring other concerns.

import (
	"context"
	"io"
	"log"
	"os"
	"path"

	utils "web-service/src/utils"

	guuid "github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	// Change this when we deploy to Amazon S3
	endpoint = "localhost:9000"
	// This is of course not secure, TODO store this in a
	// configuration file, not under VCS.
	accessKeyID     = "minioadmin"
	secretAccessKey = "minioadmin"
	// TODO experiment with =true, it failed previously but I think
	// it was because of another problem (agevorgyan)
	useSSL = false

	bucketName = "user-files"

	loggingPath = "logging/wrapMinIO.log"
)

var (
	minioClient *minio.Client

	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
	DebugLogger   *log.Logger
	rootCtx       context.Context
)

func InitializeS3Support() {
	var err error
	if err = os.Mkdir("logging", os.ModePerm); err != nil {
		ErrorLogger.Fatal(err)
	}
	diagFile, err := utils.CreateFileIfNotExists(loggingPath)
	if err != nil {
		log.Fatalln("not able to initialize logger")
	}

	// It is easy to make a shadowing mistake here
	WarningLogger = utils.GetLoggerPkgScoped("WARNING: ", diagFile)
	ErrorLogger = utils.GetLoggerPkgScoped("ERROR: ", diagFile)
	DebugLogger = utils.GetLoggerPkgScoped("DEBUG: ", diagFile)
	rootCtx = context.Background()

	minioClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		ErrorLogger.Fatalln(err)
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
			DebugLogger.Printf("bucket %#v is already on the server\n", bucketName)
		} else {
			ErrorLogger.Fatal(err)
		}
	} else {
		DebugLogger.Println("created new bucket")
	}
}

func StoreFileByUUID(id guuid.UUID, file io.Reader, fileName string) error {
	objectName := path.Join(id.String(), "/", fileName)
	// - "application/octet-stream" means binary file, which we have
	// internally in Go here (io.Reader)
	// - "-1" means "unknown size"
	_, err := minioClient.PutObject(rootCtx, bucketName, objectName, file, -1, minio.PutObjectOptions{ContentType: "application/text", ContentEncoding: "utf-8"})

	if err != nil {
		ErrorLogger.Println(err)
	} else {
		DebugLogger.Println("PutObject success")
	}
	return err
}

func LoadFileByUUID(id guuid.UUID, fileName string) io.Reader {
	objectName := path.Join(id.String(), "/", fileName)
	v, err := minioClient.GetObject(rootCtx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		ErrorLogger.Println(err)
		return nil
	} else {
		return v
	}
}
