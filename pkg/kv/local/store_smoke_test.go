package local_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/dubin555/azlake/pkg/kv"
	_ "github.com/dubin555/azlake/pkg/kv/local"
	"github.com/dubin555/azlake/pkg/kv/kvparams"
)

func TestLocalKV(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "localkv-test-*")
	if err != nil {
		t.Fatalf("create temp dir: %s", err)
	}
	defer os.RemoveAll(tmpDir)

	ctx := context.Background()
	store, err := kv.Open(ctx, kvparams.Config{
		Type: "local",
		Local: &kvparams.Local{
			Path:         tmpDir,
			PrefetchSize: 10,
		},
	})
	if err != nil {
		t.Fatalf("open store: %s", err)
	}
	defer store.Close()

	partitionKey := []byte("test-partition")
	key := []byte("greeting")
	value := []byte("hello world")

	// Set
	if err := store.Set(ctx, partitionKey, key, value); err != nil {
		t.Fatalf("Set: %s", err)
	}

	// Get
	res, err := store.Get(ctx, partitionKey, key)
	if err != nil {
		t.Fatalf("Get: %s", err)
	}
	if string(res.Value) != string(value) {
		t.Fatalf("value mismatch: got %q, want %q", res.Value, value)
	}

	// Delete
	if err := store.Delete(ctx, partitionKey, key); err != nil {
		t.Fatalf("Delete: %s", err)
	}

	// Verify not found
	_, err = store.Get(ctx, partitionKey, key)
	if !errors.Is(err, kv.ErrNotFound) {
		t.Fatalf("expected ErrNotFound after delete, got: %v", err)
	}
}
