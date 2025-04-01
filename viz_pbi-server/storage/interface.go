package storage

import (
	"context"
	"io"
)

// StorageProvider defines the interface for storage operations
type StorageProvider interface {
	// UploadFile uploads a file to the specified container
	UploadFile(ctx context.Context, containerName string, fileName string, data io.Reader) error

	// GetFile retrieves a file from the specified container
	GetFile(ctx context.Context, containerName string, fileName string) ([]byte, error)
}
