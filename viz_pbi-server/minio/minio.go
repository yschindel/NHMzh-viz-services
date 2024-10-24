package minio

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var minioClient *minio.Client

func init() {
	endpoint := getEnv("MINIO_ENDPOINT", "localhost")
	port := getEnv("MINIO_PORT", "9000")
	accessKeyID := getEnv("MINIO_ACCESS_KEY", "ROOTUSER")
	secretAccessKey := getEnv("MINIO_SECRET_KEY", "CHANGEME123")
	useSSL := getEnv("MINIO_USE_SSL", "false") == "true"

	var err error
	minioUrl := fmt.Sprintf("%s:%s", endpoint, port)
	minioClient, err = minio.New(minioUrl, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func GetFile(bucketName string, objectName string) ([]byte, error) {
	log.Printf("Fetching file: '%s' from bucket: '%s'", objectName, bucketName)
	obj, err := minioClient.GetObject(context.Background(), bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer obj.Close()

	stat, err := obj.Stat()
	if err != nil {
		return nil, err
	}

	buf := make([]byte, stat.Size)
	_, err = obj.Read(buf)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return buf, nil
}
