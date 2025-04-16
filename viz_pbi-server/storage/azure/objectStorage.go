package azure

import (
	"context"
	"fmt"
	"io"
	"strings"

	"viz_pbi-server/logger"
	"viz_pbi-server/models"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/google/uuid"
)

type BlobStorage struct {
	client *azblob.Client
	logger *logger.Logger
	config *Config
}

// NewBlobStorage creates a new Azure Blob Storage instance
func NewBlobStorage(config *Config) (*BlobStorage, error) {
	client, err := azblob.NewClientFromConnectionString(config.ConnectionString, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %v", err)
	}

	return &BlobStorage{
		client: client,
		logger: logger.WithFields(logger.Fields{"component": "storage/azure/objectStorage.go"}),
		config: config,
	}, nil
}

// CreateContainerIfNotExists creates a container if it doesn't exist
func (b *BlobStorage) CreateContainerIfNotExists(ctx context.Context, containerName string) error {
	b.logger.WithFields(logger.Fields{"container": containerName}).Debug("Checking if container exists")

	_, err := b.client.CreateContainer(ctx, containerName, nil)
	if err != nil {
		if strings.Contains(err.Error(), "ContainerAlreadyExists") {
			b.logger.WithFields(logger.Fields{"container": containerName}).Debug("Container already exists")
			return nil
		}
		b.logger.WithFields(logger.Fields{
			"container": containerName,
			"error":     err,
		}).Error("Failed to create container")
		return fmt.Errorf("failed to create container: %v", err)
	}

	b.logger.WithFields(logger.Fields{"container": containerName}).Debug("Container created successfully")
	return nil
}

// UploadBlob uploads a file to Azure Blob Storage
func (b *BlobStorage) UploadBlob(ctx context.Context, blobData models.BlobData) (models.BlobData, error) {
	b.logger.WithFields(logger.Fields{
		"container":   blobData.Container,
		"filename":    blobData.Filename,
		"projectName": blobData.Project,
		"timestamp":   blobData.Timestamp,
	}).Debug("Starting file upload")

	if err := b.CreateContainerIfNotExists(ctx, blobData.Container); err != nil {
		return blobData, err
	}

	if blobData.BlobID == "" {
		blobData.BlobID = uuid.New().String()
	}

	// metadata for the blob
	metadata := map[string]*string{
		"projectName": &blobData.Project,
		"fileName":    &blobData.Filename,
		"timestamp":   &blobData.Timestamp,
	}

	// readable metadata for logging
	readableMetadata := map[string]string{
		"projectName": blobData.Project,
		"fileName":    blobData.Filename,
		"timestamp":   blobData.Timestamp,
	}

	_, err := b.client.UploadStream(ctx, blobData.Container, blobData.BlobID, blobData.Blob, &azblob.UploadStreamOptions{
		Metadata: metadata,
	})
	if err != nil {
		b.logger.WithFields(logger.Fields{
			"container": blobData.Container,
			"blobID":    blobData.BlobID,
			"metadata":  readableMetadata,
			"error":     err,
		}).Error("Failed to upload file")
		return blobData, fmt.Errorf("failed to upload file: %v", err)
	}

	b.logger.WithFields(logger.Fields{
		"container": blobData.Container,
		"blobID":    blobData.BlobID,
		"metadata":  readableMetadata,
	}).Info("File uploaded successfully")

	return blobData, nil
}

// GetBlob gets a file from Azure Blob Storage
func (b *BlobStorage) GetBlob(ctx context.Context, containerName string, fileName string) ([]byte, error) {
	b.logger.WithFields(logger.Fields{
		"container": containerName,
		"file":      fileName,
	}).Debug("Starting file download")

	response, err := b.client.DownloadStream(ctx, containerName, fileName, nil)
	if err != nil {
		b.logger.WithFields(logger.Fields{
			"container": containerName,
			"file":      fileName,
			"error":     err,
		}).Error("Failed to download file")
		return nil, fmt.Errorf("failed to download file: %v", err)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		b.logger.WithFields(logger.Fields{
			"container": containerName,
			"file":      fileName,
			"error":     err,
		}).Error("Failed to read file data")
		return nil, fmt.Errorf("failed to read file data: %v", err)
	}

	b.logger.WithFields(logger.Fields{
		"container": containerName,
		"file":      fileName,
	}).Info("File downloaded successfully")
	return data, nil
}

// GetBlobMetadata gets the metadata from a file in Azure Blob Storage
func (b *BlobStorage) GetBlobMetadata(ctx context.Context, containerName string, fileName string) (map[string]*string, error) {
	blobClient := b.client.ServiceClient().NewContainerClient(containerName).NewBlobClient(fileName)
	response, err := blobClient.GetProperties(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get blob metadata: %v", err)
	}

	return response.Metadata, nil
}

// Container returns the default container name
func (s *BlobStorage) Container() string {
	return s.config.Container
}
