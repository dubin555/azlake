package samplerepo

import (
	"context"

	"github.com/dubin555/azlake/pkg/catalog"
)

func PopulateSampleRepo(ctx context.Context, repo *catalog.Repository, cat *catalog.Catalog, pathProvider interface{}, adapter interface{}, user interface{}) error {
	return nil // stub
}

func AddBranchProtection(ctx context.Context, repo *catalog.Repository, cat *catalog.Catalog) error {
	return nil // stub
}
