package block

import (
	"context"
	"strconv"
)

const (
	MetadataBlockstoreTypeKey = "blockstore_type"
	MetadataBlockstoreCount   = "blockstore_count"
)

// MetadataProvider is a metadata provider that reports a single blockstore type.
type MetadataProvider struct {
	blockstoreType  string
	blockstoreCount int
}

// GetMetadata returns metadata with a blockstore type(s).
// in case there is more than one blockstore - the count is also reported.
func (p *MetadataProvider) GetMetadata(_ context.Context) (map[string]string, error) {
	m := map[string]string{
		MetadataBlockstoreTypeKey: p.blockstoreType,
	}
	if p.blockstoreCount > 1 {
		m[MetadataBlockstoreCount] = strconv.Itoa(p.blockstoreCount)
	}
	return m, nil
}

// NewMetadataProvider creates a MetadataProvider for the given blockstore type.
func NewMetadataProvider(blockstoreType string) *MetadataProvider {
	return &MetadataProvider{
		blockstoreType:  blockstoreType,
		blockstoreCount: 1,
	}
}
