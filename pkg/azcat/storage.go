package azcat

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// ObjectStorage abstracts object content storage (local FS, Azure Blob, etc.)
type ObjectStorage interface {
	// Put stores object content and returns physical address, size, and checksum
	Put(ctx context.Context, repo, key string, reader io.Reader) (physicalAddr string, size int64, checksum string, err error)
	// Get retrieves object content
	Get(ctx context.Context, repo, key string) (io.ReadCloser, error)
	// Delete removes object content
	Delete(ctx context.Context, repo, key string) error
}

// LocalStorage stores objects on the local filesystem
type LocalStorage struct {
	baseDir string
}

func NewLocalStorage(baseDir string) *LocalStorage {
	_ = os.MkdirAll(baseDir, 0755)
	return &LocalStorage{baseDir: baseDir}
}

func (s *LocalStorage) Put(_ context.Context, repo, key string, reader io.Reader) (string, int64, string, error) {
	physPath := filepath.Join(s.baseDir, repo, key)
	if err := os.MkdirAll(filepath.Dir(physPath), 0755); err != nil {
		return "", 0, "", err
	}
	f, err := os.Create(physPath)
	if err != nil {
		return "", 0, "", err
	}
	h := sha256.New()
	size, err := io.Copy(io.MultiWriter(f, h), reader)
	f.Close()
	if err != nil {
		return "", 0, "", err
	}
	checksum := fmt.Sprintf("%x", h.Sum(nil))
	return physPath, size, checksum, nil
}

func (s *LocalStorage) Get(_ context.Context, repo, key string) (io.ReadCloser, error) {
	physPath := filepath.Join(s.baseDir, repo, key)
	return os.Open(physPath)
}

func (s *LocalStorage) Delete(_ context.Context, repo, key string) error {
	physPath := filepath.Join(s.baseDir, repo, key)
	return os.Remove(physPath)
}
