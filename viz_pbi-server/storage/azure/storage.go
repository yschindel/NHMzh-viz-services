package azure

import (
	"context"
	"io"
	"net/url"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

type BlobStorage struct {
	serviceURL azblob.ServiceURL
}

func NewBlobStorage(config *Config) (*BlobStorage, error) {
	credential, err := azblob.NewSharedKeyCredential(config.AccountName, config.AccountKey)
	if err != nil {
		return nil, err
	}

	pipeline := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	serviceURL, err := url.Parse(config.EndpointURL)
	if err != nil {
		return nil, err
	}

	return &BlobStorage{
		serviceURL: azblob.NewServiceURL(*serviceURL, pipeline),
	}, nil
}

// UploadFile uploads a file to Azure Blob Storage
func (b *BlobStorage) UploadFile(ctx context.Context, containerName string, fileName string, data io.Reader) error {
	containerURL := b.serviceURL.NewContainerURL(containerName)
	blobURL := containerURL.NewBlockBlobURL(fileName)

	_, err := azblob.UploadStreamToBlockBlob(ctx, data, blobURL, azblob.UploadStreamToBlockBlobOptions{})
	return err
}

// GetFile gets a file from Azure Blob Storage
func (b *BlobStorage) GetFile(ctx context.Context, containerName string, fileName string) ([]byte, error) {
	containerURL := b.serviceURL.NewContainerURL(containerName)
	blobURL := containerURL.NewBlockBlobURL(fileName)

	response, err := blobURL.Download(ctx, 0, azblob.CountToEnd, azblob.BlobAccessConditions{}, false, azblob.ClientProvidedKeyOptions{})
	if err != nil {
		return nil, err
	}

	bodyStream := response.Body(azblob.RetryReaderOptions{})
	defer bodyStream.Close()

	return io.ReadAll(bodyStream)
}
