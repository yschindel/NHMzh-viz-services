package storage

import (
	"context"
	"viz_pbi-server/models"
)

// StorageProvider defines the interface for storage operations
type StorageProvider interface {
	// Container returns the default container name
	Container() string

	// UploadBlob uploads a file to the specified container
	UploadBlob(ctx context.Context, blobData models.BlobData) (string, error)

	// GetFile retrieves a file from the specified container
	GetBlob(ctx context.Context, containerName string, fileName string) ([]byte, error)
}
