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
	var err error
	minioClient, err = minio.New("localhost:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("ROOTUSER", "CHANGEME123", ""),
		Secure: false,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
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
