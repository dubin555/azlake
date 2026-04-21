package factory

import (
	"context"
	"fmt"

	"github.com/dubin555/azlake/pkg/block"
	"github.com/dubin555/azlake/pkg/block/azure"
	"github.com/dubin555/azlake/pkg/block/local"
	"github.com/dubin555/azlake/pkg/block/params"
)

// BuildBlockAdapter returns a block adapter by type.
// Currently only "azure" and "local" are supported.
func BuildBlockAdapter(ctx context.Context, blockstoreType string, azureParams *params.Azure, localPath string) (block.Adapter, error) {
	switch blockstoreType {
	case block.BlockstoreTypeAzure:
		if azureParams == nil {
			return nil, fmt.Errorf("azure params required for azure blockstore")
		}
		return azure.NewAdapter(ctx, *azureParams)
	case block.BlockstoreTypeLocal:
		if localPath == "" {
			return nil, fmt.Errorf("local path required for local blockstore")
		}
		return local.NewAdapter(localPath)
	default:
		return nil, fmt.Errorf("unsupported blockstore type: %s (supported: azure, local)", blockstoreType)
	}
}
