package factory

import (
	"context"
	"fmt"
	"strings"

	"github.com/dubin555/azlake/pkg/block"
	"github.com/dubin555/azlake/pkg/block/azure"
	"github.com/dubin555/azlake/pkg/block/local"
	"github.com/dubin555/azlake/pkg/block/params"
	"github.com/dubin555/azlake/pkg/config"
	"github.com/dubin555/azlake/pkg/logging"
)

// BuildBlockAdapter creates a block adapter from config.
// Supports "azure" and "local" blockstore types.
func BuildBlockAdapter(ctx context.Context, _ interface{}, c config.AdapterConfig) (block.Adapter, error) {
	blockstore := strings.ToLower(c.BlockstoreType())
	logging.FromContext(ctx).
		WithField("type", blockstore).
		Info("initialize blockstore adapter")
	switch blockstore {
	case block.BlockstoreTypeLocal:
		p, err := c.BlockstoreLocalParams()
		if err != nil {
			return nil, err
		}
		adapter, err := local.NewAdapter(p.Path,
			local.WithAllowedExternalPrefixes(p.AllowedExternalPrefixes),
			local.WithImportEnabled(p.ImportEnabled),
		)
		if err != nil {
			return nil, fmt.Errorf("local adapter with path %s: %w", p.Path, err)
		}
		return adapter, nil
	case block.BlockstoreTypeAzure:
		p, err := c.BlockstoreAzureParams()
		if err != nil {
			return nil, err
		}
		return azure.NewAdapter(ctx, p)
	default:
		return nil, fmt.Errorf("%w '%s' (supported: azure, local)",
			block.ErrInvalidAddress, blockstore)
	}
}

// BuildBlockAdapterFromParams creates a block adapter from raw params (no config needed).
func BuildBlockAdapterFromParams(ctx context.Context, blockstoreType string, azureParams *params.Azure, localPath string) (block.Adapter, error) {
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
