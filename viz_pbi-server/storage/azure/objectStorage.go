package azure

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"

	"viz_pbi-server/logger"
	"viz_pbi-server/models"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/google/uuid"
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
		logger:     logger.WithFields(logger.Fields{"component": "storage/azure/objectStorage.go"}),
		config:     config,
	}, nil
}

// CreateContainerIfNotExists creates a container if it doesn't exist
func (b *BlobStorage) CreateContainerIfNotExists(ctx context.Context, containerName string) error {
	b.logger.WithFields(logger.Fields{"container": containerName}).Debug("Checking if container exists")

	containerURL := b.serviceURL.NewContainerURL(containerName)
	_, err := containerURL.Create(ctx, azblob.Metadata{}, azblob.PublicAccessNone)
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

// UploadFile uploads a file to Azure Blob Storage
func (b *BlobStorage) UploadBlob(ctx context.Context, blobData models.BlobData) (models.BlobData, error) {
	b.logger.WithFields(logger.Fields{
		"container":   blobData.Container,
		"filename":    blobData.Filename,
		"projectName": blobData.Project,
	}).Debug("Starting file upload")

	// Create container if it doesn't exist
	if err := b.CreateContainerIfNotExists(ctx, blobData.Container); err != nil {
		return blobData, err
	}
	containerURL := b.serviceURL.NewContainerURL(blobData.Container)

	// add the service URL to the blob data so it can be used in the sql writer
	blobData.StorageServiceURL = b.serviceURL.String()

	// generate a guid for the blob if it doesn't exist
	if blobData.BlobID == "" {
		blobData.BlobID = uuid.New().String()
	}

	blobMetadata := azblob.Metadata{
		"projectName": blobData.Project,
		"fileName":    blobData.Filename,
		"timestamp":   blobData.Timestamp,
	}

	blobURL := containerURL.NewBlockBlobURL(blobData.BlobID)

	_, err := azblob.UploadStreamToBlockBlob(ctx, blobData.Blob, blobURL, azblob.UploadStreamToBlockBlobOptions{
		Metadata: blobMetadata,
	})
	if err != nil {
		b.logger.WithFields(logger.Fields{
			"container": blobData.Container,
			"blobID":    blobData.BlobID,
			"metadata":  blobMetadata,
			"error":     err,
		}).Error("Failed to upload file")
		return blobData, fmt.Errorf("failed to upload file: %v", err)
	}

	b.logger.WithFields(logger.Fields{
		"container": blobData.Container,
		"blobID":    blobData.BlobID,
		"metadata":  blobMetadata,
	}).Info("File uploaded successfully")

	return blobData, nil
}

// GetFile gets a file from Azure Blob Storage
func (b *BlobStorage) GetBlob(ctx context.Context, containerName string, fileName string) ([]byte, error) {
	b.logger.WithFields(logger.Fields{
		"container": containerName,
		"file":      fileName,
	}).Debug("Starting file download")

	containerURL := b.serviceURL.NewContainerURL(containerName)
	blobURL := containerURL.NewBlockBlobURL(fileName)

	response, err := blobURL.Download(ctx, 0, azblob.CountToEnd, azblob.BlobAccessConditions{}, false, azblob.ClientProvidedKeyOptions{})
	if err != nil {
		b.logger.WithFields(logger.Fields{
			"container": containerName,
			"file":      fileName,
			"error":     err,
		}).Error("Failed to download file")
		return nil, fmt.Errorf("failed to download file: %v", err)
	}

	bodyStream := response.Body(azblob.RetryReaderOptions{})
	defer bodyStream.Close()

	data, err := io.ReadAll(bodyStream)
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

// Container returns the default container name
func (s *BlobStorage) Container() string {
	return s.config.Container
}
