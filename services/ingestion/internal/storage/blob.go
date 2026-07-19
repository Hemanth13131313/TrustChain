package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type BlobRepository interface {
	Save(ctx context.Context, objectName string, reader io.Reader) (string, error)
}

// LocalBlobRepository is an MVP implementation storing blobs on the local disk
type LocalBlobRepository struct {
	baseDir string
}

func NewLocalBlobRepository(baseDir string) (*LocalBlobRepository, error) {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("creating blob storage dir: %w", err)
	}
	return &LocalBlobRepository{baseDir: baseDir}, nil
}

// Save writes the contents of reader to a file and returns its storage reference
func (r *LocalBlobRepository) Save(ctx context.Context, objectName string, reader io.Reader) (string, error) {
	path := filepath.Join(r.baseDir, objectName)
	
	// Create subdirectories if necessary
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return "", fmt.Errorf("creating object subdirs: %w", err)
	}

	out, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("creating object file: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, reader); err != nil {
		return "", fmt.Errorf("writing object data: %w", err)
	}

	// For local storage, the ref is just the file path
	return path, nil
}
