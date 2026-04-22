package azcat

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

// CosmosKVStore implements the same interface as KVStore but backed by Azure CosmosDB.
// Uses a single container with partition key = first path segment of the key (e.g., "repos").
// Documents have: id (full key, url-encoded), partitionKey, value (JSON-encoded bytes).
type CosmosKVStore struct {
	client *azcosmos.ContainerClient
}

type cosmosDoc struct {
	ID           string          `json:"id"`
	Key          string          `json:"key"` // original unencoded key for querying
	PartitionKey string          `json:"pk"`
	Value        json.RawMessage `json:"value"`
}

// OpenCosmosKV creates a CosmosDB-backed KV store.
// If key is empty, uses DefaultAzureCredential.
func OpenCosmosKV(endpoint, key, database, container string) (*CosmosKVStore, error) {
	var client *azcosmos.Client
	var err error

	if key != "" {
		cred, credErr := azcosmos.NewKeyCredential(key)
		if credErr != nil {
			return nil, fmt.Errorf("cosmos key credential: %w", credErr)
		}
		client, err = azcosmos.NewClientWithKey(endpoint, cred, nil)
	} else {
		cred, credErr := azidentity.NewDefaultAzureCredential(nil)
		if credErr != nil {
			return nil, fmt.Errorf("cosmos default credential: %w", credErr)
		}
		client, err = azcosmos.NewClient(endpoint, cred, nil)
	}
	if err != nil {
		return nil, fmt.Errorf("cosmos client: %w", err)
	}

	containerClient, err := client.NewContainer(database, container)
	if err != nil {
		return nil, fmt.Errorf("cosmos container client: %w", err)
	}

	return &CosmosKVStore{client: containerClient}, nil
}

// partitionKey extracts the first path segment as partition key.
// e.g., "repos/my-repo/branches/main" → "repos"
func partitionKey(key string) string {
	if idx := strings.Index(key, "/"); idx > 0 {
		return key[:idx]
	}
	return key
}

// encodeID replaces '/' with '|' for use as a CosmosDB document ID.
// CosmosDB IDs cannot contain '/' characters — even URL-encoded %2F gets decoded.
func encodeID(key string) string {
	return strings.ReplaceAll(key, "/", "|")
}

// decodeID reverses encodeID.
func decodeID(id string) string {
	return strings.ReplaceAll(id, "|", "/")
}

func (s *CosmosKVStore) Close() error {
	// CosmosDB client doesn't need explicit close
	return nil
}

func (s *CosmosKVStore) Get(key string) ([]byte, error) {
	ctx := context.Background()
	pk := azcosmos.NewPartitionKeyString(partitionKey(key))

	resp, err := s.client.ReadItem(ctx, pk, encodeID(key), nil)
	if err != nil {
		var respErr *azcore.ResponseError
		if errors.As(err, &respErr) && respErr.StatusCode == http.StatusNotFound {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("cosmos get %q: %w", key, err)
	}

	var doc cosmosDoc
	if err := json.Unmarshal(resp.Value, &doc); err != nil {
		return nil, fmt.Errorf("cosmos unmarshal %q: %w", key, err)
	}
	return doc.Value, nil
}

func (s *CosmosKVStore) Set(key string, value []byte) error {
	ctx := context.Background()
	pk := azcosmos.NewPartitionKeyString(partitionKey(key))

	doc := cosmosDoc{
		ID:           encodeID(key),
		Key:          key,
		PartitionKey: partitionKey(key),
		Value:        json.RawMessage(value),
	}
	data, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("cosmos marshal: %w", err)
	}

	_, err = s.client.UpsertItem(ctx, pk, data, nil)
	if err != nil {
		return fmt.Errorf("cosmos set %q: %w", key, err)
	}
	return nil
}

func (s *CosmosKVStore) Delete(key string) error {
	ctx := context.Background()
	pk := azcosmos.NewPartitionKeyString(partitionKey(key))

	_, err := s.client.DeleteItem(ctx, pk, encodeID(key), nil)
	if err != nil {
		var respErr *azcore.ResponseError
		if errors.As(err, &respErr) && respErr.StatusCode == http.StatusNotFound {
			return nil // already deleted
		}
		return fmt.Errorf("cosmos delete %q: %w", key, err)
	}
	return nil
}

func (s *CosmosKVStore) Scan(prefix string, after string, limit int, fn func(key string, value []byte) error) error {
	ctx := context.Background()
	pk := azcosmos.NewPartitionKeyString(partitionKey(prefix))

	query := "SELECT * FROM c WHERE STARTSWITH(c.key, @prefix)"
	params := []azcosmos.QueryParameter{
		{Name: "@prefix", Value: prefix},
	}
	if after != "" {
		query += " AND c.key > @after"
		params = append(params, azcosmos.QueryParameter{Name: "@after", Value: after})
	}
	query += " ORDER BY c.key"

	queryOpts := &azcosmos.QueryOptions{
		QueryParameters: params,
	}

	pager := s.client.NewQueryItemsPager(query, pk, queryOpts)
	count := 0
	for pager.More() {
		if limit > 0 && count >= limit {
			break
		}
		page, err := pager.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("cosmos scan: %w", err)
		}
		for _, item := range page.Items {
			if limit > 0 && count >= limit {
				break
			}
			var doc cosmosDoc
			if err := json.Unmarshal(item, &doc); err != nil {
				continue
			}
			if err := fn(doc.Key, doc.Value); err != nil {
				return err
			}
			count++
		}
	}
	return nil
}

func (s *CosmosKVStore) DeletePrefix(prefix string) error {
	ctx := context.Background()
	pk := azcosmos.NewPartitionKeyString(partitionKey(prefix))

	query := "SELECT c.id FROM c WHERE STARTSWITH(c.key, @prefix)"
	queryOpts := &azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@prefix", Value: prefix},
		},
	}

	pager := s.client.NewQueryItemsPager(query, pk, queryOpts)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("cosmos delete prefix scan: %w", err)
		}
		for _, item := range page.Items {
			var doc struct {
				ID string `json:"id"`
			}
			if err := json.Unmarshal(item, &doc); err != nil {
				continue
			}
			_, _ = s.client.DeleteItem(ctx, pk, doc.ID, nil)
		}
	}
	return nil
}

// JSON helpers — same interface as KVStore
func (s *CosmosKVStore) GetJSON(key string, v interface{}) error {
	data, err := s.Get(key)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

func (s *CosmosKVStore) SetJSON(key string, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return s.Set(key, data)
}

// Ensure CosmosKVStore can be used wherever KVStore is used
// Note: Both implement the same methods but Go doesn't have structural typing for concrete types.
// The Catalog accepts *KVStore. To use CosmosDB, we need an interface. See KV interface below.

// KV is the interface that both KVStore (BadgerDB) and CosmosKVStore implement
type KV interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
	Delete(key string) error
	Scan(prefix string, after string, limit int, fn func(key string, value []byte) error) error
	DeletePrefix(prefix string) error
	GetJSON(key string, v interface{}) error
	SetJSON(key string, v interface{}) error
	Close() error
}
