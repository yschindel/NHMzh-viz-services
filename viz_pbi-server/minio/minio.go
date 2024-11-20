package minio

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

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
		if minioErr, ok := err.(minio.ErrorResponse); ok && minioErr.Code == "NoSuchKey" {
			log.Printf("File '%s' not found in bucket '%s'", objectName, bucketName)
			return nil, fmt.Errorf("file '%s' not found in bucket '%s'", objectName, bucketName)
		}
		log.Printf("Failed to get file '%s': %v", objectName, err)
		return nil, err
	}
	defer obj.Close()

	stat, err := obj.Stat()
	if err != nil {
		log.Printf("Failed to get file '%s' stats: %v", objectName, err)
		return nil, err
	}

	buf := make([]byte, stat.Size)
	_, err = obj.Read(buf)
	if err != nil && err != io.EOF {
		log.Printf("Failed to read file '%s': %v", objectName, err)
		return nil, err
	}

	log.Printf("File '%s' fetched successfully", objectName)
	return buf, nil
}

func ListFiles(bucketName string, prefix string) ([]string, error) {
	// make sure prefix ends with a /
	if !strings.HasSuffix(prefix, "/") {
		prefix = prefix + "/"
	}
	log.Printf("Listing files in bucket: '%s' with prefix: '%s'", bucketName, prefix)
	objectCh := minioClient.ListObjects(context.Background(), bucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	files := []string{}
	for object := range objectCh {
		files = append(files, object.Key)
	}

	return files, nil
}

func ListProjects(bucketName string) ([]string, error) {
	objects := minioClient.ListObjects(context.Background(), bucketName, minio.ListObjectsOptions{
		Prefix:    "",
		Recursive: false,
	})

	projects := []string{}
	for object := range objects {
		projects = append(projects, strings.Split(object.Key, "/")[0])
	}

	return projects, nil
}

func ListAllFiles(bucketName string) ([]string, error) {
	objects := minioClient.ListObjects(context.Background(), bucketName, minio.ListObjectsOptions{
		Prefix:    "",
		Recursive: true,
	})

	files := []string{}
	for object := range objects {
		files = append(files, object.Key)
	}

	return files, nil
}
