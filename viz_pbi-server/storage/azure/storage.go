package azure

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"

	"viz_pbi-server/logger"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

type BlobStorage struct {
	serviceURL azblob.ServiceURL
	logger     *logger.Logger
	config     *Config
}

// NewBlobStorage creates a new Azure Blob Storage instance
func NewBlobStorage(config *Config) (*BlobStorage, error) {
	// Create a default request pipeline using your storage account name and account key
	credential, err := azblob.NewSharedKeyCredential(config.AccountName, config.AccountKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create shared key credential: %v", err)
	}

	pipeline := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	// Create a service URL
	serviceURL, err := url.Parse(config.EndpointURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse endpoint URL: %v", err)
	}

	return &BlobStorage{
		serviceURL: azblob.NewServiceURL(*serviceURL, pipeline),
		logger:     logger.WithFields(logger.Fields{"component": "azure_storage"}),
		config:     config,
	}, nil
}

// CreateContainerIfNotExists creates a container if it doesn't exist
func (b *BlobStorage) CreateContainerIfNotExists(ctx context.Context, containerName string) error {
	b.logger.Debug("Creating container if not exists: %s", containerName)

	containerURL := b.serviceURL.NewContainerURL(containerName)
	_, err := containerURL.Create(ctx, azblob.Metadata{}, azblob.PublicAccessNone)
	if err != nil && !strings.Contains(err.Error(), "ContainerAlreadyExists") {
		b.logger.Error("Failed to create container %s: %v", containerName, err)
		return fmt.Errorf("failed to create container: %v", err)
	}

	b.logger.Debug("Container exists or was created: %s", containerName)
	return nil
}

// UploadFile uploads a file to Azure Blob Storage
func (b *BlobStorage) UploadFile(ctx context.Context, containerName string, fileName string, data io.Reader) error {
	b.logger.Debug("Starting file upload: container=%s, file=%s", containerName, fileName)

	// Create container if it doesn't exist
	if err := b.CreateContainerIfNotExists(ctx, containerName); err != nil {
		return err
	}

	containerURL := b.serviceURL.NewContainerURL(containerName)
	blobURL := containerURL.NewBlockBlobURL(fileName)

	_, err := azblob.UploadStreamToBlockBlob(ctx, data, blobURL, azblob.UploadStreamToBlockBlobOptions{})
	if err != nil {
		b.logger.Error("Failed to upload file: container=%s, file=%s, error=%v", containerName, fileName, err)
		return fmt.Errorf("failed to upload file: %v", err)
	}

	b.logger.Info("File uploaded successfully: container=%s, file=%s", containerName, fileName)
	return nil
}

// GetFile gets a file from Azure Blob Storage
func (b *BlobStorage) GetFile(ctx context.Context, containerName string, fileName string) ([]byte, error) {
	b.logger.Debug("Starting file download: container=%s, file=%s", containerName, fileName)

	containerURL := b.serviceURL.NewContainerURL(containerName)
	blobURL := containerURL.NewBlockBlobURL(fileName)

	response, err := blobURL.Download(ctx, 0, azblob.CountToEnd, azblob.BlobAccessConditions{}, false, azblob.ClientProvidedKeyOptions{})
	if err != nil {
		b.logger.Error("Failed to download file: container=%s, file=%s, error=%v", containerName, fileName, err)
		return nil, fmt.Errorf("failed to download file: %v", err)
	}

	bodyStream := response.Body(azblob.RetryReaderOptions{})
	defer bodyStream.Close()

	data, err := io.ReadAll(bodyStream)
	if err != nil {
		b.logger.Error("Failed to read file data: container=%s, file=%s, error=%v", containerName, fileName, err)
		return nil, fmt.Errorf("failed to read file data: %v", err)
	}

	b.logger.Info("File downloaded successfully: container=%s, file=%s", containerName, fileName)
	return data, nil
}

// Container returns the default container name
func (s *BlobStorage) Container() string {
	return s.config.Container
}
