package azcat

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/sas"
)

// AzureBlobStorage stores objects in Azure Blob Storage
type AzureBlobStorage struct {
	client        *azblob.Client
	containerName string
	accountName   string
	sharedKeyCred *azblob.SharedKeyCredential // nil if using DefaultAzureCredential
}

// NewAzureBlobStorage creates an Azure Blob Storage backend.
// If accountKey is empty, it falls back to DefaultAzureCredential (az login, managed identity, etc.)
func NewAzureBlobStorage(accountName, accountKey, containerName string) (*AzureBlobStorage, error) {
	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)

	var client *azblob.Client
	var sharedKeyCred *azblob.SharedKeyCredential
	var err error

	if accountKey != "" {
		cred, credErr := azblob.NewSharedKeyCredential(accountName, accountKey)
		if credErr != nil {
			return nil, fmt.Errorf("azure shared key credential: %w", credErr)
		}
		sharedKeyCred = cred
		client, err = azblob.NewClientWithSharedKeyCredential(serviceURL, cred, nil)
	} else {
		cred, credErr := azidentity.NewDefaultAzureCredential(nil)
		if credErr != nil {
			return nil, fmt.Errorf("azure default credential: %w", credErr)
		}
		client, err = azblob.NewClient(serviceURL, cred, nil)
	}
	if err != nil {
		return nil, fmt.Errorf("azure blob client: %w", err)
	}

	return &AzureBlobStorage{
		client:        client,
		containerName: containerName,
		accountName:   accountName,
		sharedKeyCred: sharedKeyCred,
	}, nil
}

func (s *AzureBlobStorage) blobName(repo, key string) string {
	return repo + "/" + key
}

func (s *AzureBlobStorage) Put(ctx context.Context, repo, key string, reader io.Reader) (string, int64, string, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return "", 0, "", fmt.Errorf("reading content: %w", err)
	}

	h := sha256.New()
	h.Write(data)
	checksum := fmt.Sprintf("%x", h.Sum(nil))
	size := int64(len(data))

	blobName := s.blobName(repo, key)
	_, err = s.client.UploadBuffer(ctx, s.containerName, blobName, data, &azblob.UploadBufferOptions{})
	if err != nil {
		return "", 0, "", fmt.Errorf("uploading to azure blob %q: %w", blobName, err)
	}

	physAddr := fmt.Sprintf("az://%s/%s", s.containerName, blobName)
	return physAddr, size, checksum, nil
}

func (s *AzureBlobStorage) Get(ctx context.Context, repo, key string) (io.ReadCloser, error) {
	blobName := s.blobName(repo, key)
	resp, err := s.client.DownloadStream(ctx, s.containerName, blobName, nil)
	if err != nil {
		return nil, fmt.Errorf("downloading azure blob %q: %w", blobName, err)
	}
	return resp.Body, nil
}

func (s *AzureBlobStorage) Delete(ctx context.Context, repo, key string) error {
	blobName := s.blobName(repo, key)
	_, err := s.client.DeleteBlob(ctx, s.containerName, blobName, nil)
	if err != nil {
		return fmt.Errorf("deleting azure blob %q: %w", blobName, err)
	}
	return nil
}

// GetSASURL generates a read-only SAS URL for the given object, valid for the specified duration.
// Returns empty string if SAS generation is not supported (e.g., using DefaultAzureCredential).
func (s *AzureBlobStorage) GetSASURL(repo, key string, expiry time.Duration) (string, error) {
	if s.sharedKeyCred == nil {
		return "", fmt.Errorf("SAS URL generation requires shared key credential")
	}

	blobName := s.blobName(repo, key)
	now := time.Now().UTC()

	sasQueryParams, err := sas.BlobSignatureValues{
		Protocol:      sas.ProtocolHTTPS,
		StartTime:     now.Add(-5 * time.Minute), // clock skew tolerance
		ExpiryTime:    now.Add(expiry),
		Permissions:   (&sas.BlobPermissions{Read: true}).String(),
		ContainerName: s.containerName,
		BlobName:      blobName,
	}.SignWithSharedKey(s.sharedKeyCred)
	if err != nil {
		return "", fmt.Errorf("generating SAS: %w", err)
	}

	// Build the full URL
	sasURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s?%s",
		s.accountName, s.containerName, blobName, sasQueryParams.Encode())
	return sasURL, nil
}

// Ensure interface compliance
var _ ObjectStorage = (*AzureBlobStorage)(nil)

// SASCapable is an optional interface that storage backends can implement
// to support generating pre-signed/SAS URLs for direct browser access.
type SASCapable interface {
	GetSASURL(repo, key string, expiry time.Duration) (string, error)
}

// Compile-time check
var _ SASCapable = (*AzureBlobStorage)(nil)
