# azlake Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Extract an Azure-native data version control system from lakeFS into a standalone open-source project.

**Architecture:** Copy-and-trim approach — copy needed packages from lakeFS source at `$LAKEFS_SRC/`, rename module to `github.com/dubin555/azlake`, delete non-Azure code, simplify auth to API-key-only, and build a slim React WebUI with DuckDB WASM. Five phases following dependency order: Foundation → Engine → Business → Interface → WebUI.

**Tech Stack:** Go 1.25, Azure SDK for Go (azblob, azcosmos, azidentity), BadgerDB, Pebble, Protobuf, chi router, cobra CLI, React 18 + Vite + TypeScript + DuckDB WASM.

**Source codebase:** `$LAKEFS_SRC/` (referred to as `$SRC` throughout)

---

## Phase 1: Foundation

### Task 1: Initialize project and Go module

**Files:**
- Create: `go.mod`
- Create: `go.sum` (generated)
- Create: `Makefile`
- Create: `README.md`

- [ ] **Step 1: Create project directory and init Go module**

```bash
mkdir -p $PROJECT_ROOT
cd $PROJECT_ROOT
git init
go mod init github.com/dubin555/azlake
```

- [ ] **Step 2: Create Makefile**

```makefile
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOFMT=$(GOCMD)fmt

AZLAKE_BINARY=azlake
AZLAKECTL_BINARY=azlakectl

.PHONY: all build test fmt clean

all: build

build: build-server build-cli

build-server:
	$(GOBUILD) -o $(AZLAKE_BINARY) ./cmd/azlake

build-cli:
	$(GOBUILD) -o $(AZLAKECTL_BINARY) ./cmd/azlakectl

test:
	$(GOTEST) -count=1 -race ./...

fmt:
	$(GOFMT) -w .

clean:
	rm -f $(AZLAKE_BINARY) $(AZLAKECTL_BINARY)
```

- [ ] **Step 3: Create README.md**

```markdown
# azlake

Azure-native data version control. Git for your data lake.

Extracted from [lakeFS](https://github.com/treeverse/lakefs), stripped to Azure-only.

## Quick Start (local mode)

\`\`\`bash
go build -o azlake ./cmd/azlake
./azlake run --config azlake.local.yaml
\`\`\`

## License

Apache 2.0
```

- [ ] **Step 4: Commit**

```bash
git add go.mod Makefile README.md
git commit -m "feat: initialize azlake project"
```

---

### Task 2: Copy KV store interface and supporting packages

This task copies the foundational packages that have no internal dependencies: logging, ident, KV interface.

**Files:**
- Create: `pkg/logging/` (copy from `$SRC/pkg/logging/`)
- Create: `pkg/ident/` (copy from `$SRC/pkg/ident/`)
- Create: `pkg/kv/store.go` (copy from `$SRC/pkg/kv/store.go`)
- Create: `pkg/kv/kvparams/database.go` (copy from `$SRC/pkg/kv/kvparams/database.go`)

- [ ] **Step 1: Copy logging package**

```bash
mkdir -p pkg/logging
cp $SRC/pkg/logging/*.go pkg/logging/
# Remove test files
rm -f pkg/logging/*_test.go
```

- [ ] **Step 2: Rename imports in logging**

```bash
find pkg/logging -name '*.go' -exec sed -i '' 's|github.com/treeverse/lakefs|github.com/dubin555/azlake|g' {} +
```

- [ ] **Step 3: Copy ident package**

```bash
mkdir -p pkg/ident
cp $SRC/pkg/ident/*.go pkg/ident/
rm -f pkg/ident/*_test.go
find pkg/ident -name '*.go' -exec sed -i '' 's|github.com/treeverse/lakefs|github.com/dubin555/azlake|g' {} +
```

- [ ] **Step 4: Copy KV store interface**

```bash
mkdir -p pkg/kv
cp $SRC/pkg/kv/store.go pkg/kv/
cp $SRC/pkg/kv/store_message.go pkg/kv/ 2>/dev/null || true
cp $SRC/pkg/kv/secondary_index.go pkg/kv/ 2>/dev/null || true
cp $SRC/pkg/kv/errors.go pkg/kv/ 2>/dev/null || true

# Copy any additional non-test, non-mock top-level kv files
for f in $SRC/pkg/kv/*.go; do
  fname=$(basename "$f")
  if [[ "$fname" != *_test.go ]] && [[ "$fname" != *mock* ]]; then
    cp "$f" pkg/kv/
  fi
done

# Copy kvparams
mkdir -p pkg/kv/kvparams
cp $SRC/pkg/kv/kvparams/*.go pkg/kv/kvparams/

# Rename imports
find pkg/kv -name '*.go' -exec sed -i '' 's|github.com/treeverse/lakefs|github.com/dubin555/azlake|g' {} +
```

- [ ] **Step 5: Run go mod tidy to pull dependencies**

```bash
go mod tidy
```

- [ ] **Step 6: Verify it compiles**

Run: `go build ./pkg/kv/... ./pkg/logging/... ./pkg/ident/...`
Expected: Clean build, no errors.

- [ ] **Step 7: Commit**

```bash
git add pkg/logging pkg/ident pkg/kv go.mod go.sum
git commit -m "feat: add KV store interface, logging, and ident packages"
```

---

### Task 3: Add KV drivers — CosmosDB and Local (BadgerDB)

**Files:**
- Create: `pkg/kv/cosmosdb/` (copy from `$SRC/pkg/kv/cosmosdb/`)
- Create: `pkg/kv/local/` (copy from `$SRC/pkg/kv/local/`)

- [ ] **Step 1: Copy CosmosDB driver**

```bash
mkdir -p pkg/kv/cosmosdb
for f in $SRC/pkg/kv/cosmosdb/*.go; do
  fname=$(basename "$f")
  if [[ "$fname" != *_test.go ]]; then
    cp "$f" pkg/kv/cosmosdb/
  fi
done
find pkg/kv/cosmosdb -name '*.go' -exec sed -i '' 's|github.com/treeverse/lakefs|github.com/dubin555/azlake|g' {} +
```

- [ ] **Step 2: Copy Local (BadgerDB) driver**

```bash
mkdir -p pkg/kv/local
for f in $SRC/pkg/kv/local/*.go; do
  fname=$(basename "$f")
  if [[ "$fname" != *_test.go ]]; then
    cp "$f" pkg/kv/local/
  fi
done
find pkg/kv/local -name '*.go' -exec sed -i '' 's|github.com/treeverse/lakefs|github.com/dubin555/azlake|g' {} +
```

- [ ] **Step 3: Remove registration of other drivers**

Check `pkg/kv/store.go` or any init files for references to dynamodb, postgres, mem drivers. Remove those imports/registrations if present. Only cosmosdb and local should be registered.

- [ ] **Step 4: Run go mod tidy and verify**

```bash
go mod tidy
go build ./pkg/kv/...
```

Expected: Clean build. `go.sum` should include Azure SDK and BadgerDB dependencies.

- [ ] **Step 5: Write a basic test for Local KV**

Create `pkg/kv/local/store_smoke_test.go`:

```go
package local_test

import (
	"context"
	"os"
	"testing"

	"github.com/dubin555/azlake/pkg/kv"
	kvlocal "github.com/dubin555/azlake/pkg/kv/local"
	"github.com/dubin555/azlake/pkg/kv/kvparams"
)

func TestLocalKVGetSetDelete(t *testing.T) {
	dir, err := os.MkdirTemp("", "azlake-kv-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	cfg := kvparams.Config{
		Type: "local",
		Local: &kvparams.Local{
			Path:          dir,
			PrefetchSize:  256,
			EnableLogging: false,
		},
	}
	store, err := (&kvlocal.Driver{}).Open(context.Background(), cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	ctx := context.Background()
	pk := []byte("test-partition")
	key := []byte("test-key")
	val := []byte("test-value")

	// Set
	if err := store.Set(ctx, pk, key, val); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Get
	result, err := store.Get(ctx, pk, key)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if string(result.Value) != string(val) {
		t.Fatalf("expected %q, got %q", val, result.Value)
	}

	// Delete
	if err := store.Delete(ctx, pk, key); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Get after delete should return ErrNotFound
	_, err = store.Get(ctx, pk, key)
	if err != kv.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got: %v", err)
	}
}
```

- [ ] **Step 6: Run test**

Run: `go test ./pkg/kv/local/ -v -run TestLocalKVGetSetDelete`
Expected: PASS

- [ ] **Step 7: Commit**

```bash
git add pkg/kv/cosmosdb pkg/kv/local go.mod go.sum
git commit -m "feat: add CosmosDB and Local (BadgerDB) KV drivers"
```

---

### Task 4: Copy Block adapter interface + Azure + Local implementations

**Files:**
- Create: `pkg/block/` core files (copy from `$SRC/pkg/block/`)
- Create: `pkg/block/params/` (copy from `$SRC/pkg/block/params/`)
- Create: `pkg/block/azure/` (copy from `$SRC/pkg/block/azure/`)
- Create: `pkg/block/local/` (copy from `$SRC/pkg/block/local/`)
- Create: `pkg/block/factory/build.go` (simplified)
- Create: `pkg/api/apiutil/` (copy needed utility types)

- [ ] **Step 1: Copy block core interface files**

```bash
mkdir -p pkg/block
for f in $SRC/pkg/block/*.go; do
  fname=$(basename "$f")
  if [[ "$fname" != *_test.go ]] && [[ "$fname" != *mock* ]]; then
    cp "$f" pkg/block/
  fi
done
mkdir -p pkg/block/params
cp $SRC/pkg/block/params/*.go pkg/block/params/
find pkg/block -maxdepth 2 -name '*.go' -exec sed -i '' 's|github.com/treeverse/lakefs|github.com/dubin555/azlake|g' {} +
```

- [ ] **Step 2: Copy Azure adapter**

```bash
mkdir -p pkg/block/azure
for f in $SRC/pkg/block/azure/*.go; do
  fname=$(basename "$f")
  if [[ "$fname" != *_test.go ]]; then
    cp "$f" pkg/block/azure/
  fi
done
find pkg/block/azure -name '*.go' -exec sed -i '' 's|github.com/treeverse/lakefs|github.com/dubin555/azlake|g' {} +
```

- [ ] **Step 3: Copy Local adapter**

```bash
mkdir -p pkg/block/local
for f in $SRC/pkg/block/local/*.go; do
  fname=$(basename "$f")
  if [[ "$fname" != *_test.go ]]; then
    cp "$f" pkg/block/local/
  fi
done
find pkg/block/local -name '*.go' -exec sed -i '' 's|github.com/treeverse/lakefs|github.com/dubin555/azlake|g' {} +
```

- [ ] **Step 4: Copy apiutil if needed by block**

Check if any block files import `pkg/api/apiutil`. If so:

```bash
mkdir -p pkg/api/apiutil
for f in $SRC/pkg/api/apiutil/*.go; do
  fname=$(basename "$f")
  if [[ "$fname" != *_test.go ]]; then
    cp "$f" pkg/api/apiutil/
  fi
done
find pkg/api/apiutil -name '*.go' -exec sed -i '' 's|github.com/treeverse/lakefs|github.com/dubin555/azlake|g' {} +
```

- [ ] **Step 5: Create simplified factory**

Create `pkg/block/factory/build.go`:

```go
package factory

import (
	"context"
	"fmt"
	"strings"

	"github.com/dubin555/azlake/pkg/block"
	"github.com/dubin555/azlake/pkg/block/azure"
	"github.com/dubin555/azlake/pkg/block/local"
	"github.com/dubin555/azlake/pkg/config"
	"github.com/dubin555/azlake/pkg/logging"
)

func BuildBlockAdapter(ctx context.Context, c config.AdapterConfig) (block.Adapter, error) {
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
		return nil, fmt.Errorf("%w '%s' please choose one of %s",
			block.ErrInvalidAddress, blockstore,
			[]string{block.BlockstoreTypeLocal, block.BlockstoreTypeAzure})
	}
}
```

Note: This may not compile yet because `config.AdapterConfig` doesn't exist. That's OK — we'll create the config package in the next task. For now just create the file.

- [ ] **Step 6: Run go mod tidy and attempt build**

```bash
go mod tidy
go build ./pkg/block/... 2>&1 | head -20
```

Expected: May have errors related to missing `config` or `stats` packages. Note them — we'll fix in subsequent tasks. The core `pkg/block/azure/` and `pkg/block/local/` should compile independently.

- [ ] **Step 7: Commit**

```bash
git add pkg/block pkg/api/apiutil go.mod go.sum
git commit -m "feat: add block storage adapters (Azure Blob + Local FS)"
```

---

### Task 5: Copy config and supporting packages

**Files:**
- Create: `pkg/config/` (copy from `$SRC/pkg/config/`, simplified)
- Create: `pkg/stats/` (copy minimal stats interface if needed)
- Create: `pkg/httputil/` (copy from `$SRC/pkg/httputil/`)
- Create: `pkg/validator/` (copy from `$SRC/pkg/validator/` if it exists)

- [ ] **Step 1: Copy config package**

```bash
mkdir -p pkg/config
for f in $SRC/pkg/config/*.go; do
  fname=$(basename "$f")
  if [[ "$fname" != *_test.go ]]; then
    cp "$f" pkg/config/
  fi
done
find pkg/config -name '*.go' -exec sed -i '' 's|github.com/treeverse/lakefs|github.com/dubin555/azlake|g' {} +
```

- [ ] **Step 2: Remove non-Azure config options**

In the config files, look for S3, GCS, DynamoDB configuration structs and references. Comment them out or delete them. Keep only:
- Azure Blob storage config
- CosmosDB config  
- Local config (both KV and Block)
- Listen address, auth secret, logging

- [ ] **Step 3: Copy other supporting packages as needed**

Check compile errors and copy any missing packages that block compilation. Common ones:

```bash
# Copy httputil if needed
if grep -r "pkg/httputil" pkg/ 2>/dev/null | grep -q .; then
  mkdir -p pkg/httputil
  for f in $SRC/pkg/httputil/*.go; do
    fname=$(basename "$f")
    if [[ "$fname" != *_test.go ]]; then
      cp "$f" pkg/httputil/
    fi
  done
  find pkg/httputil -name '*.go' -exec sed -i '' 's|github.com/treeverse/lakefs|github.com/dubin555/azlake|g' {} +
fi

# Copy stats interface if needed (often just an interface + noop implementation)
if grep -r "pkg/stats" pkg/ 2>/dev/null | grep -q .; then
  mkdir -p pkg/stats
  for f in $SRC/pkg/stats/*.go; do
    fname=$(basename "$f")
    if [[ "$fname" != *_test.go ]] && [[ "$fname" != *mock* ]]; then
      cp "$f" pkg/stats/
    fi
  done
  find pkg/stats -name '*.go' -exec sed -i '' 's|github.com/treeverse/lakefs|github.com/dubin555/azlake|g' {} +
fi
```

- [ ] **Step 4: Iteratively fix compile errors**

```bash
go mod tidy
go build ./pkg/... 2>&1 | head -30
```

For each missing package, either copy it from `$SRC` or create a minimal stub. The goal is `go build ./pkg/kv/... ./pkg/block/... ./pkg/config/... ./pkg/logging/...` passes.

- [ ] **Step 5: Commit**

```bash
git add pkg/ go.mod go.sum
git commit -m "feat: add config, httputil, stats, and supporting packages"
```

---

### Task 6: Phase 1 compilation checkpoint

- [ ] **Step 1: Verify all foundation packages compile**

```bash
go build ./pkg/kv/...
go build ./pkg/block/...
go build ./pkg/config/...
go build ./pkg/logging/...
go build ./pkg/ident/...
```

Expected: All clean builds.

- [ ] **Step 2: Run any existing tests**

```bash
go test ./pkg/kv/local/ -v -count=1
```

Expected: PASS

- [ ] **Step 3: Tag milestone**

```bash
git tag phase1-foundation
```

---

## Phase 2: Core Engine

### Task 7: Copy Graveler version engine

This is the largest copy — ~13K lines. The graveler has subdirectories for committed (MetaRange), ref (branch/commit/tag management), staging, and settings.

**Files:**
- Create: `pkg/graveler/` (entire directory tree from `$SRC/pkg/graveler/`)

- [ ] **Step 1: Copy entire graveler directory**

```bash
mkdir -p pkg/graveler
# Copy all .go and .proto files, preserving directory structure
cd $SRC && find pkg/graveler -name '*.go' -not -name '*_test.go' -not -path '*/mock/*' | while read f; do
  mkdir -p $PROJECT_ROOT/$(dirname "$f")
  cp "$f" $PROJECT_ROOT/"$f"
done
cd $PROJECT_ROOT

# Copy proto files
cd $SRC && find pkg/graveler -name '*.proto' | while read f; do
  mkdir -p $PROJECT_ROOT/$(dirname "$f")
  cp "$f" $PROJECT_ROOT/"$f"
done
cd $PROJECT_ROOT

# Copy generated .pb.go files
cd $SRC && find pkg/graveler -name '*.pb.go' | while read f; do
  mkdir -p $PROJECT_ROOT/$(dirname "$f")
  cp "$f" $PROJECT_ROOT/"$f"
done
cd $PROJECT_ROOT
```

- [ ] **Step 2: Rename all imports**

```bash
find pkg/graveler -name '*.go' -exec sed -i '' 's|github.com/treeverse/lakefs|github.com/dubin555/azlake|g' {} +
```

- [ ] **Step 3: Remove PullRequest code**

```bash
# Remove PR-related files
find pkg/graveler -name '*pull*' -o -name '*pr_*' | xargs rm -f 2>/dev/null

# In graveler.proto, the PullRequestData and PullRequestStatus messages can stay
# (removing proto fields breaks .pb.go — safer to leave them unused)
```

In `pkg/graveler/graveler.go`, find and remove (or comment out) any PullRequest-related methods. Search for `PullRequest`, `CreatePullRequest`, `UpdatePullRequest`, `ListPullRequests`, `GetPullRequest`.

- [ ] **Step 4: Remove branch protection references**

In `pkg/graveler/graveler.go`, the `protectedBranchesManager` is used in Commit and Reset. For the simplest approach, replace the protection checks with no-ops:

Find lines like:
```go
isProtected, err := g.protectedBranchesManager.IsBlocked(...)
if isProtected { return "", ErrCommitToProtectedBranch }
```

Replace with:
```go
// Branch protection removed in azlake
```

- [ ] **Step 5: Copy missing dependencies referenced by graveler**

Graveler depends on several packages. Copy them:

```bash
# permissions package (used for permission constants)
if grep -r "pkg/permissions" pkg/graveler/ 2>/dev/null | grep -q .; then
  mkdir -p pkg/permissions
  for f in $SRC/pkg/permissions/*.go; do
    fname=$(basename "$f")
    if [[ "$fname" != *_test.go ]]; then
      cp "$f" pkg/permissions/
    fi
  done
  find pkg/permissions -name '*.go' -exec sed -i '' 's|github.com/treeverse/lakefs|github.com/dubin555/azlake|g' {} +
fi

# validator package
if grep -r "pkg/validator" pkg/graveler/ 2>/dev/null | grep -q .; then
  mkdir -p pkg/validator
  cp $SRC/pkg/validator/*.go pkg/validator/ 2>/dev/null || true
  rm -f pkg/validator/*_test.go
  find pkg/validator -name '*.go' -exec sed -i '' 's|github.com/treeverse/lakefs|github.com/dubin555/azlake|g' {} +
fi
```

- [ ] **Step 6: Iteratively fix compile errors**

```bash
go mod tidy
go build ./pkg/graveler/... 2>&1 | head -30
```

This will likely have errors from missing packages or removed features. For each:
- Missing package → copy from `$SRC` or create minimal stub
- Reference to removed feature (PullRequest, branch protection) → delete the code
- Reference to hooks → stub out (we'll add webhook hooks in Phase 3)

Repeat until `go build ./pkg/graveler/...` passes.

- [ ] **Step 7: Commit**

```bash
git add pkg/ go.mod go.sum
git commit -m "feat: add graveler version engine"
```

---

### Task 8: Copy upload utilities and verify Phase 2

**Files:**
- Create: `pkg/upload/` (copy from `$SRC/pkg/upload/`)

- [ ] **Step 1: Copy upload package**

```bash
mkdir -p pkg/upload
for f in $SRC/pkg/upload/*.go; do
  fname=$(basename "$f")
  if [[ "$fname" != *_test.go ]]; then
    cp "$f" pkg/upload/
  fi
done
find pkg/upload -name '*.go' -exec sed -i '' 's|github.com/treeverse/lakefs|github.com/dubin555/azlake|g' {} +
```

- [ ] **Step 2: Verify full Phase 2 build**

```bash
go mod tidy
go build ./pkg/...
```

Expected: Clean build of all packages so far.

- [ ] **Step 3: Commit and tag**

```bash
git add pkg/upload go.mod go.sum
git commit -m "feat: add upload utilities"
git tag phase2-engine
```

---

## Phase 3: Business Layer

### Task 9: Copy catalog package

**Files:**
- Create: `pkg/catalog/` (copy from `$SRC/pkg/catalog/`)

- [ ] **Step 1: Copy catalog**

```bash
mkdir -p pkg/catalog
for f in $SRC/pkg/catalog/*.go; do
  fname=$(basename "$f")
  if [[ "$fname" != *_test.go ]] && [[ "$fname" != *mock* ]]; then
    cp "$f" pkg/catalog/
  fi
done
find pkg/catalog -name '*.go' -exec sed -i '' 's|github.com/treeverse/lakefs|github.com/dubin555/azlake|g' {} +
```

- [ ] **Step 2: Remove PullRequest references from catalog**

Search for PullRequest-related methods in catalog files and remove them:

```bash
grep -n -i "pullrequest\|pull_request\|PullRequest" pkg/catalog/*.go
```

Remove those methods. They should be self-contained (not called by other retained code).

- [ ] **Step 3: Fix compile errors**

```bash
go mod tidy
go build ./pkg/catalog/... 2>&1 | head -30
```

Fix missing dependencies by copying from `$SRC` or stubbing.

- [ ] **Step 4: Commit**

```bash
git add pkg/catalog go.mod go.sum
git commit -m "feat: add catalog coordination layer"
```

---

### Task 10: Create simplified API Key auth

Instead of copying lakeFS's 12K-line auth module, create a simplified version that only handles API key pairs.

**Files:**
- Create: `pkg/auth/model/model.go`
- Create: `pkg/auth/model/model.proto`
- Create: `pkg/auth/model/model.pb.go`
- Create: `pkg/auth/keys/keys.go`
- Create: `pkg/auth/crypt/encryption.go`
- Create: `pkg/auth/service.go`
- Create: `pkg/auth/middleware.go`
- Create: `pkg/auth/context.go`
- Test: `pkg/auth/service_test.go`

- [ ] **Step 1: Copy auth model (user and credential structs)**

```bash
mkdir -p pkg/auth/model
cp $SRC/pkg/auth/model/model.go pkg/auth/model/
cp $SRC/pkg/auth/model/model.pb.go pkg/auth/model/
cp $SRC/pkg/auth/model/model.proto pkg/auth/model/ 2>/dev/null || true
cp $SRC/pkg/auth/model/validation.go pkg/auth/model/ 2>/dev/null || true
find pkg/auth/model -name '*.go' -exec sed -i '' 's|github.com/treeverse/lakefs|github.com/dubin555/azlake|g' {} +
```

- [ ] **Step 2: Copy encryption and key generation**

```bash
mkdir -p pkg/auth/crypt
cp $SRC/pkg/auth/crypt/*.go pkg/auth/crypt/
rm -f pkg/auth/crypt/*_test.go
find pkg/auth/crypt -name '*.go' -exec sed -i '' 's|github.com/treeverse/lakefs|github.com/dubin555/azlake|g' {} +

mkdir -p pkg/auth/keys
cp $SRC/pkg/auth/keys/*.go pkg/auth/keys/
rm -f pkg/auth/keys/*_test.go
find pkg/auth/keys -name '*.go' -exec sed -i '' 's|github.com/treeverse/lakefs|github.com/dubin555/azlake|g' {} +
```

- [ ] **Step 3: Copy context helpers and params**

```bash
cp $SRC/pkg/auth/context.go pkg/auth/
cp $SRC/pkg/auth/errors.go pkg/auth/
cp $SRC/pkg/auth/token.go pkg/auth/
mkdir -p pkg/auth/params
cp $SRC/pkg/auth/params/*.go pkg/auth/params/
find pkg/auth -maxdepth 1 -name '*.go' -exec sed -i '' 's|github.com/treeverse/lakefs|github.com/dubin555/azlake|g' {} +
find pkg/auth/params -name '*.go' -exec sed -i '' 's|github.com/treeverse/lakefs|github.com/dubin555/azlake|g' {} +
```

- [ ] **Step 4: Create simplified auth service**

Create `pkg/auth/service.go` — a stripped version that only manages API key credentials:

```go
package auth

// This is a simplified auth service for azlake.
// It only supports API key authentication (access_key_id + secret_access_key).
// No RBAC, no groups, no policies, no ACLs.
//
// The original lakeFS auth service is at $SRC/pkg/auth/service.go (1,331 lines).
// This version should be ~200-300 lines.
//
// Key operations:
// - CreateCredentials(username) → (accessKeyID, secretAccessKey)
// - GetCredentials(accessKeyID) → (*Credential, error)
// - GetUser(username) → (*User, error)
// - CreateUser(username) → error
// - ListUsers() → ([]*User, error)
```

Copy the relevant parts from `$SRC/pkg/auth/service.go` and `$SRC/pkg/auth/basic_service.go`, keeping only:
- `CreateCredentials`
- `GetCredentials` / `GetCredentialsForUser`
- `AddCredentials` / `DeleteCredentials`
- `CreateUser` / `GetUser` / `ListUsers` / `DeleteUser`

Remove all methods related to: groups, policies, ACLs, external auth, OIDC, LDAP.

- [ ] **Step 5: Copy request_auth.go (API key extraction from HTTP request)**

```bash
cp $SRC/pkg/auth/request_auth.go pkg/auth/
sed -i '' 's|github.com/treeverse/lakefs|github.com/dubin555/azlake|g' pkg/auth/request_auth.go
```

- [ ] **Step 6: Write auth service test**

Create `pkg/auth/service_test.go`:

```go
package auth_test

import (
	"context"
	"testing"

	// Import the auth service once it compiles
	// This test verifies basic credential lifecycle
)

func TestCreateAndGetCredentials(t *testing.T) {
	// TODO: instantiate auth service with in-memory or local KV
	// Create a user
	// Create credentials for the user
	// Get credentials by access key ID
	// Verify the secret matches
	t.Skip("implement after auth service compiles")
}
```

- [ ] **Step 7: Fix compile errors and commit**

```bash
go mod tidy
go build ./pkg/auth/... 2>&1 | head -30
# Fix errors iteratively
git add pkg/auth go.mod go.sum
git commit -m "feat: add simplified API Key auth service"
```

---

### Task 11: Create webhook hooks (extracted from lakeFS actions)

**Files:**
- Create: `pkg/hooks/webhook.go` (extracted from `$SRC/pkg/actions/webhook.go`)
- Create: `pkg/hooks/hook.go` (extracted from `$SRC/pkg/actions/hook.go`)
- Create: `pkg/hooks/action.go` (extracted from `$SRC/pkg/actions/action.go`)
- Create: `pkg/hooks/service.go` (simplified from `$SRC/pkg/actions/service.go`)

- [ ] **Step 1: Copy core hook files**

```bash
mkdir -p pkg/hooks
cp $SRC/pkg/actions/webhook.go pkg/hooks/
cp $SRC/pkg/actions/hook.go pkg/hooks/
cp $SRC/pkg/actions/action.go pkg/hooks/

# Change package name from actions to hooks
find pkg/hooks -name '*.go' -exec sed -i '' 's|package actions|package hooks|g' {} +
find pkg/hooks -name '*.go' -exec sed -i '' 's|github.com/treeverse/lakefs|github.com/dubin555/azlake|g' {} +
```

- [ ] **Step 2: Create simplified service**

Copy `$SRC/pkg/actions/service.go` to `pkg/hooks/service.go`, then strip it down:

```bash
cp $SRC/pkg/actions/service.go pkg/hooks/
sed -i '' 's|package actions|package hooks|g' pkg/hooks/service.go
sed -i '' 's|github.com/treeverse/lakefs|github.com/dubin555/azlake|g' pkg/hooks/service.go
```

Then edit `pkg/hooks/service.go` to:
- Remove all references to Lua hooks and Airflow hooks
- Remove action run history storage (the `StoreRun` / `GetRun` / `ListRuns` methods)
- Keep only the webhook dispatch logic
- In the `NewHook` factory function or equivalent, only support `type: webhook`

- [ ] **Step 3: Remove Lua and Airflow references**

```bash
# Delete any copied lua/airflow files
rm -f pkg/hooks/lua*.go pkg/hooks/airflow*.go
```

In `pkg/hooks/action.go`, find the HookType constants and remove lua/airflow:

```go
// Keep only:
const (
    HookTypeWebhook HookType = "webhook"
)
```

- [ ] **Step 4: Fix compile errors**

```bash
go mod tidy
go build ./pkg/hooks/... 2>&1 | head -30
```

Remove any remaining references to Lua VM, Airflow, or action run storage.

- [ ] **Step 5: Commit**

```bash
git add pkg/hooks go.mod go.sum
git commit -m "feat: add webhook-only hooks system"
```

---

### Task 12: Phase 3 compilation checkpoint

- [ ] **Step 1: Verify full build**

```bash
go build ./pkg/...
```

Expected: Clean build of all packages.

- [ ] **Step 2: Run any tests**

```bash
go test ./pkg/kv/local/ -v -count=1
```

- [ ] **Step 3: Tag milestone**

```bash
git tag phase3-business
```

---

## Phase 4: Interface Layer

### Task 13: Copy and trim OpenAPI spec + generated code

**Files:**
- Create: `api/swagger.yml` (trimmed from `$SRC/api/swagger.yml`)
- Create: `pkg/api/apigen/` (regenerated or copied + trimmed)

- [ ] **Step 1: Copy OpenAPI spec**

```bash
mkdir -p api
cp $SRC/api/swagger.yml api/
```

- [ ] **Step 2: Trim OpenAPI spec**

Remove endpoints from `api/swagger.yml` related to:
- All `/auth/` endpoints (users, groups, policies, ACLs) — keep only login/credentials
- All `/repositories/{repository}/pulls/` endpoints
- All `/statistics/` endpoints  
- S3 gateway endpoints

Keep:
- `/repositories` CRUD
- `/repositories/{repository}/branches` CRUD
- `/repositories/{repository}/refs/{ref}/objects` (get, stat, upload, delete)
- `/repositories/{repository}/branches/{branch}/commits` 
- `/repositories/{repository}/refs/{ref}/diff`
- `/repositories/{repository}/branches/{branch}/merge`
- `/repositories/{repository}/branches/{branch}/revert`
- `/repositories/{repository}/tags` CRUD
- `/config` endpoints
- Add: `/repositories/{repository}/refs/{ref}/objects/presign` for DuckDB

- [ ] **Step 3: Copy or regenerate API code**

Option A (simpler): Copy the generated file and trim unused types/handlers:
```bash
mkdir -p pkg/api/apigen
cp $SRC/pkg/api/apigen/lakefs.gen.go pkg/api/apigen/
cp $SRC/pkg/api/apigen/doc.go pkg/api/apigen/ 2>/dev/null || true
find pkg/api/apigen -name '*.go' -exec sed -i '' 's|github.com/treeverse/lakefs|github.com/dubin555/azlake|g' {} +
```

Option B (cleaner but requires Docker): Regenerate from the trimmed spec:
```bash
# Requires oapi-codegen
go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.5.6 \
  -package apigen \
  -generate "types,client,chi-server,spec" \
  -o pkg/api/apigen/azlake.gen.go \
  api/swagger.yml
```

Choose Option A for now (faster), clean up later.

- [ ] **Step 4: Commit**

```bash
git add api/ pkg/api/apigen go.mod go.sum
git commit -m "feat: add trimmed OpenAPI spec and generated API layer"
```

---

### Task 14: Copy and trim API controller + middleware

**Files:**
- Create: `pkg/api/controller.go` (copy from `$SRC/pkg/api/controller.go`, trim)
- Create: `pkg/api/serve.go` (copy from `$SRC/pkg/api/serve.go`)
- Create: `pkg/api/auth_middleware.go` (copy from `$SRC/pkg/api/auth_middleware.go`)
- Create: other middleware files

- [ ] **Step 1: Copy all API handler files**

```bash
for f in $SRC/pkg/api/*.go; do
  fname=$(basename "$f")
  if [[ "$fname" != *_test.go ]] && [[ "$fname" != *mock* ]] && [[ "$fname" != *gen* ]]; then
    cp "$f" pkg/api/
  fi
done
find pkg/api -maxdepth 1 -name '*.go' -exec sed -i '' 's|github.com/treeverse/lakefs|github.com/dubin555/azlake|g' {} +
```

- [ ] **Step 2: Trim controller.go**

`controller.go` is 6,495 lines. Remove methods for:
- All auth management handlers (CreateUser, ListUsers, GetGroup, CreatePolicy, etc.)
- PullRequest handlers
- Statistics handlers
- S3 gateway handlers
- Any handler referencing removed features

Keep all handlers for: repositories, branches, commits, tags, objects, diff, merge, revert.

Search patterns to find removable methods:
```bash
grep -n "func.*Controller.*Pull\|func.*Controller.*Group\|func.*Controller.*Policy\|func.*Controller.*Statistics" pkg/api/controller.go
```

- [ ] **Step 3: Update serve.go**

Remove route registrations for deleted endpoints. Remove references to S3 gateway, removed auth endpoints.

- [ ] **Step 4: Simplify auth_middleware.go**

Replace the RBAC permission checks with simple API key validation. The middleware should:
1. Extract API key from request (Basic auth header or query param)
2. Look up the key in the auth service
3. Set the authenticated user in context
4. No permission checks (all authenticated users can do everything)

- [ ] **Step 5: Fix compile errors iteratively**

```bash
go mod tidy
go build ./pkg/api/... 2>&1 | head -30
```

This will be the longest step. Work through errors one by one:
- Missing interface methods → add stubs or copy implementations
- References to removed types → delete the code
- Missing packages → copy from `$SRC`

- [ ] **Step 6: Commit**

```bash
git add pkg/api go.mod go.sum
git commit -m "feat: add REST API controller and middleware"
```

---

### Task 15: Create server binary

**Files:**
- Create: `cmd/azlake/main.go`
- Create: `cmd/azlake/cmd/root.go`
- Create: `cmd/azlake/cmd/run.go`
- Create: `cmd/azlake/cmd/setup.go`
- Create: `azlake.local.yaml`

- [ ] **Step 1: Copy and simplify server command**

```bash
mkdir -p cmd/azlake/cmd
cp $SRC/cmd/lakefs/main.go cmd/azlake/
cp $SRC/cmd/lakefs/cmd/root.go cmd/azlake/cmd/
cp $SRC/cmd/lakefs/cmd/run.go cmd/azlake/cmd/
cp $SRC/cmd/lakefs/cmd/setup.go cmd/azlake/cmd/
cp $SRC/cmd/lakefs/cmd/common_helpers.go cmd/azlake/cmd/

# Rename imports and binary name
find cmd/azlake -name '*.go' -exec sed -i '' 's|github.com/treeverse/lakefs|github.com/dubin555/azlake|g' {} +
find cmd/azlake -name '*.go' -exec sed -i '' 's|lakefs|azlake|g' {} +
```

- [ ] **Step 2: Remove unused commands from run.go**

In `cmd/azlake/cmd/run.go`:
- Remove S3 gateway startup
- Remove WebUI embedding (for now — we'll add it in Phase 5)
- Remove references to lua, actions service (replace with hooks service)
- Remove gc, migrate, flare commands

- [ ] **Step 3: Create local config file**

Create `azlake.local.yaml`:

```yaml
listen_address: "0.0.0.0:8000"

logging:
  level: DEBUG

auth:
  encrypt:
    secret_key: "THIS_IS_A_LOCAL_DEV_SECRET_CHANGE_IN_PROD"

blockstore:
  type: local
  local:
    path: /tmp/azlake/data

database:
  type: local
  local:
    path: /tmp/azlake/metadata
```

- [ ] **Step 4: Build and test server binary**

```bash
go mod tidy
go build -o azlake ./cmd/azlake
./azlake --help
```

Expected: Help output showing available commands (at minimum: `run`, `setup`).

- [ ] **Step 5: Test server starts in local mode**

```bash
mkdir -p /tmp/azlake/data /tmp/azlake/metadata
./azlake run --config azlake.local.yaml &
sleep 2
curl -s http://localhost:8000/api/v1/config | head -20
kill %1
```

Expected: Server starts and responds to health/config endpoint.

- [ ] **Step 6: Commit**

```bash
git add cmd/azlake azlake.local.yaml go.mod go.sum
git commit -m "feat: add azlake server binary with local mode"
```

---

### Task 16: Create CLI binary

**Files:**
- Create: `cmd/azlakectl/main.go`
- Create: `cmd/azlakectl/cmd/` (copy core commands from `$SRC/cmd/lakectl/cmd/`)

- [ ] **Step 1: Copy CLI framework**

```bash
mkdir -p cmd/azlakectl/cmd

# Copy main.go
cp $SRC/cmd/lakectl/main.go cmd/azlakectl/

# Copy core command files (not abuse, not auth management, not actions)
for f in root.go common.go repo.go branch.go commit.go tag.go fs.go diff.go merge.go log.go; do
  cp $SRC/cmd/lakectl/cmd/$f cmd/azlakectl/cmd/ 2>/dev/null || true
done

# Copy any utility files needed
cp $SRC/cmd/lakectl/cmd/output.go cmd/azlakectl/cmd/ 2>/dev/null || true

# Rename imports
find cmd/azlakectl -name '*.go' -exec sed -i '' 's|github.com/treeverse/lakefs|github.com/dubin555/azlake|g' {} +
find cmd/azlakectl -name '*.go' -exec sed -i '' 's|lakectl|azlakectl|g' {} +
```

- [ ] **Step 2: Remove commands for deleted features**

Remove files for:
- `abuse*.go` (load testing)
- `auth*.go` (RBAC management)
- `actions*.go` (action runs)
- `bisect*.go` (if exists)
- Any other feature-specific commands not in the core set

- [ ] **Step 3: Build CLI**

```bash
go mod tidy
go build -o azlakectl ./cmd/azlakectl
./azlakectl --help
```

Expected: Help showing core commands (repo, branch, commit, tag, fs, diff, merge, log).

- [ ] **Step 4: Test CLI against running server**

```bash
# Start server in background
./azlake run --config azlake.local.yaml &
sleep 2

# Test basic operations
./azlakectl repo list --server-address http://localhost:8000

kill %1
```

- [ ] **Step 5: Commit**

```bash
git add cmd/azlakectl go.mod go.sum
git commit -m "feat: add azlakectl CLI binary"
```

---

### Task 17: End-to-end local mode test

- [ ] **Step 1: Start server**

```bash
rm -rf /tmp/azlake
mkdir -p /tmp/azlake/data /tmp/azlake/metadata
./azlake run --config azlake.local.yaml &
sleep 3
```

- [ ] **Step 2: Setup (create initial admin credentials)**

```bash
./azlake setup --config azlake.local.yaml --user-name admin
# Note the access key ID and secret
```

- [ ] **Step 3: Test full workflow**

```bash
export AZLAKECTL_SERVER=http://localhost:8000
export AZLAKECTL_ACCESS_KEY_ID=<from setup>
export AZLAKECTL_SECRET_ACCESS_KEY=<from setup>

# Create repository
./azlakectl repo create azlake://test-repo --default-branch main --storage-namespace local:///tmp/azlake/data/test-repo

# Upload a file
echo '{"test": true}' > /tmp/test.json
./azlakectl fs upload azlake://test-repo/main/test.json --source /tmp/test.json

# List files
./azlakectl fs ls azlake://test-repo/main/

# Commit
./azlakectl commit azlake://test-repo/main -m "initial test data"

# Create branch
./azlakectl branch create azlake://test-repo/dev --source azlake://test-repo/main

# Upload to branch
echo '{"test": true, "branch": "dev"}' > /tmp/test2.json
./azlakectl fs upload azlake://test-repo/dev/test.json --source /tmp/test2.json

# Diff
./azlakectl diff azlake://test-repo/main azlake://test-repo/dev

# Commit branch
./azlakectl commit azlake://test-repo/dev -m "modify on dev branch"

# Merge
./azlakectl merge azlake://test-repo/dev azlake://test-repo/main

# Log
./azlakectl log azlake://test-repo/main

# Revert
./azlakectl branch revert azlake://test-repo/main main~1 --parent-number 1 --yes
```

- [ ] **Step 4: Verify and clean up**

```bash
kill %1  # stop server
```

Expected: All commands succeed. This validates the entire backend stack works end-to-end in local mode.

- [ ] **Step 5: Tag milestone**

```bash
git tag phase4-interface
```

---

## Phase 5: Web UI

### Task 18: Initialize slim React project

**Files:**
- Create: `webui/package.json`
- Create: `webui/vite.config.ts`
- Create: `webui/tsconfig.json`
- Create: `webui/index.html`
- Create: `webui/src/main.tsx`
- Create: `webui/src/App.tsx`

- [ ] **Step 1: Create Vite React project**

```bash
cd webui 2>/dev/null || mkdir webui
cd webui
npm create vite@latest . -- --template react-ts
npm install
npm install react-router-dom
npm install react-bootstrap bootstrap
npm install @primer/octicons-react
npm install @duckdb/duckdb-wasm apache-arrow
npm install react-diff-viewer-continued
npm install prismjs
npm install dayjs
```

- [ ] **Step 2: Configure Vite proxy for API**

Update `webui/vite.config.ts`:

```typescript
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8000',
        changeOrigin: true,
      },
    },
  },
})
```

- [ ] **Step 3: Verify dev server starts**

```bash
cd webui && npm run dev
```

Expected: Vite dev server starts on port 5173.

- [ ] **Step 4: Commit**

```bash
cd ..  # back to project root
git add webui/package.json webui/vite.config.ts webui/tsconfig.json webui/index.html webui/src/
echo "node_modules" >> webui/.gitignore
git add webui/.gitignore
git commit -m "feat: initialize slim React WebUI with Vite"
```

---

### Task 19: Add API client and basic pages

**Files:**
- Create: `webui/src/lib/api.ts` (API client)
- Create: `webui/src/pages/Repositories.tsx`
- Create: `webui/src/pages/Repository.tsx`
- Create: `webui/src/pages/ObjectBrowser.tsx`

- [ ] **Step 1: Create API client**

Extract the API client pattern from `$SRC/webui/src/lib/api/` and create a simplified version that calls azlake REST API. Key methods:

- `listRepositories()`
- `createRepository(name, storageNamespace, defaultBranch)`
- `listBranches(repo)`
- `listObjects(repo, ref, prefix)`
- `getObject(repo, ref, path)` → returns content or pre-signed URL
- `uploadObject(repo, branch, path, file)`
- `commit(repo, branch, message)`
- `diff(repo, leftRef, rightRef)`
- `merge(repo, sourceBranch, destBranch)`

- [ ] **Step 2: Create repository list page**

- [ ] **Step 3: Create object browser page**

- [ ] **Step 4: Set up React Router**

- [ ] **Step 5: Commit**

```bash
git add webui/src/
git commit -m "feat: add API client and basic repository/object browser pages"
```

---

### Task 20: Add file renderers and DuckDB WASM

**Files:**
- Create: `webui/src/renderers/duckdb.tsx` (modified from `$SRC`)
- Create: `webui/src/renderers/data.tsx` (copy from `$SRC`)
- Create: `webui/src/renderers/simple.tsx` (copy from `$SRC`)
- Create: `webui/src/renderers/editor.tsx` (copy from `$SRC`)
- Create: `webui/src/renderers/index.tsx` (copy from `$SRC`)

- [ ] **Step 1: Copy renderers from lakeFS**

```bash
mkdir -p webui/src/renderers
cp $SRC/webui/src/pages/repositories/repository/fileRenderers/*.tsx webui/src/renderers/
cp $SRC/webui/src/pages/repositories/repository/fileRenderers/*.jsx webui/src/renderers/
```

- [ ] **Step 2: Modify DuckDB to use pre-signed URLs**

Replace the S3 protocol approach in `webui/src/renderers/duckdb.tsx`:

```typescript
// Replace the S3 configuration section with pre-signed URL fetch:

const LAKEFS_URI_PATTERN = /^(['"]?)(lakefs:\/\/(.*))(['"])\s*$/;

async function resolvePresignedUrl(repo: string, ref: string, path: string): Promise<string> {
    const response = await fetch(
        `/api/v1/repositories/${encodeURIComponent(repo)}/refs/${encodeURIComponent(ref)}/objects/presign?path=${encodeURIComponent(path)}`
    );
    if (!response.ok) {
        throw new Error(`Failed to get presigned URL: ${response.statusText}`);
    }
    const data = await response.json();
    return data.url;
}

// In extractFiles, instead of mapping to s3:// URLs, resolve pre-signed HTTP URLs
async function extractFiles(conn: AsyncDuckDBConnection, sql: string): Promise<{ [name: string]: string }> {
    const tokenized = await conn.bindings.tokenize(sql);
    const fileMap: { [name: string]: string } = {};
    let prev = 0;
    tokenized.offsets.forEach((offset, i) => {
        let currentToken = sql.length;
        if (i < tokenized.offsets.length - 1) {
            currentToken = tokenized.offsets[i + 1];
        }
        const part = sql.substring(prev, currentToken);
        prev = currentToken;
        if (tokenized.types[i] === DUCKDB_STRING_CONSTANT) {
            const matches = part.match(LAKEFS_URI_PATTERN);
            if (matches !== null) {
                const uri = matches[2];
                // Parse lakefs://repo/ref/path
                const parts = uri.replace('lakefs://', '').split('/');
                const repo = parts[0];
                const ref = parts[1];
                const path = parts.slice(2).join('/');
                fileMap[uri] = `__RESOLVE__:${repo}:${ref}:${path}`;
            }
        }
    });
    return fileMap;
}

// In runDuckDBQuery, resolve the presigned URLs before registering
export async function runDuckDBQuery(sql: string): Promise<arrow.Table<any>> {
    const db = await getDuckDB();
    const conn = await db.connect();
    try {
        const fileMap = await extractFiles(conn, sql);
        const resolvedFiles: { [name: string]: string } = {};
        for (const [name, ref] of Object.entries(fileMap)) {
            if (ref.startsWith('__RESOLVE__:')) {
                const [_, repo, refId, path] = ref.split(':');
                resolvedFiles[name] = await resolvePresignedUrl(repo, refId, path);
            }
        }
        await Promise.all(
            Object.entries(resolvedFiles).map(([name, url]) =>
                db.registerFileURL(name, url, DuckDBDataProtocol.HTTP, true)
            )
        );
        const result = await conn.query(sql);
        await Promise.all(
            Object.keys(resolvedFiles).map((name) => db.dropFile(name))
        );
        return result;
    } finally {
        await conn.close();
    }
}
```

- [ ] **Step 3: Fix import paths in renderers**

Update all renderer imports to use the new project paths instead of lakeFS paths.

- [ ] **Step 4: Test DuckDB renderer loads in dev mode**

```bash
cd webui && npm run dev
```

Open browser, navigate to a file, verify the DuckDB renderer appears (even if the query fails without a running server).

- [ ] **Step 5: Commit**

```bash
git add webui/src/renderers/
git commit -m "feat: add file renderers with DuckDB WASM (pre-signed URL mode)"
```

---

### Task 21: Add diff visualization

**Files:**
- Create: `webui/src/pages/Compare.tsx`
- Create: `webui/src/components/ObjectsDiff.tsx`
- Create: `webui/src/pages/CommitHistory.tsx`

- [ ] **Step 1: Extract diff components from lakeFS**

Copy the relevant diff components:

```bash
# Find and copy diff-related components
find $SRC/webui/src -name '*diff*' -o -name '*Diff*' -o -name '*compare*' -o -name '*Compare*' | head -20
```

Extract `ObjectsDiff` component and the `CompareBranches` page, adapting imports for the new project structure.

- [ ] **Step 2: Create Compare page**

Wire up a `/compare` route that takes two refs and shows:
- File tree diff (added/removed/changed)
- Click on a file to see text diff (using `react-diff-viewer-continued`)

- [ ] **Step 3: Create commit history page**

Simple list of commits with message, author, date, and link to diff vs parent.

- [ ] **Step 4: Commit**

```bash
git add webui/src/pages/ webui/src/components/
git commit -m "feat: add diff visualization and commit history pages"
```

---

### Task 22: Add pre-signed URL API endpoint (backend)

The WebUI DuckDB integration needs a pre-signed URL endpoint that doesn't exist in lakeFS's REST API (it was handled by the S3 Gateway).

**Files:**
- Modify: `pkg/api/controller.go` (add presign handler)
- Modify: `api/swagger.yml` (add endpoint)

- [ ] **Step 1: Add presign handler to controller**

Add to `pkg/api/controller.go`:

```go
func (c *Controller) GetObjectPresignURL(w http.ResponseWriter, r *http.Request, repository, ref string, params apigen.GetObjectPresignURLParams) {
	ctx := r.Context()
	repo, err := c.Catalog.GetRepository(ctx, repository)
	if err != nil {
		writeError(w, r, http.StatusNotFound, err)
		return
	}
	entry, err := c.Catalog.GetEntry(ctx, repository, ref, params.Path, catalog.GetEntryParams{})
	if err != nil {
		writeError(w, r, http.StatusNotFound, err)
		return
	}
	objectPointer := block.ObjectPointer{
		StorageID:        repo.StorageID,
		StorageNamespace: repo.StorageNamespace,
		IdentifierType:   entry.AddressType.ToIdentifierType(),
		Identifier:       entry.PhysicalAddress,
	}
	preSignedURL, expiry, err := c.BlockAdapter.GetPreSignedURL(ctx, objectPointer, block.PreSignModeRead, "")
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"url":        preSignedURL,
		"expiry":     expiry,
	})
}
```

- [ ] **Step 2: Register the route**

Add the route in `serve.go` or wherever routes are registered.

- [ ] **Step 3: Test the endpoint**

```bash
# With server running
curl "http://localhost:8000/api/v1/repositories/test-repo/refs/main/objects/presign?path=test.json"
```

Expected: Returns a JSON with `url` field (for Azure, a SAS URL; for local, a local file path).

- [ ] **Step 4: Commit**

```bash
git add pkg/api/ api/
git commit -m "feat: add pre-signed URL endpoint for WebUI DuckDB"
```

---

### Task 23: Embed WebUI in server binary

**Files:**
- Modify: `cmd/azlake/cmd/run.go` (embed webui dist)
- Modify: `Makefile` (add webui build step)

- [ ] **Step 1: Add WebUI build to Makefile**

```makefile
UI_DIR=webui
UI_BUILD_DIR=$(UI_DIR)/dist

build-ui:
	cd $(UI_DIR) && npm ci && npm run build

build: build-ui build-server build-cli
```

- [ ] **Step 2: Embed dist in Go binary**

In `cmd/azlake/cmd/run.go`, use `embed.FS` to serve the built WebUI:

```go
//go:embed all:../../webui/dist
var webUIFS embed.FS

// In the server setup, serve the embedded files at /
```

- [ ] **Step 3: Build and test**

```bash
make build-ui
make build-server
./azlake run --config azlake.local.yaml &
sleep 2
curl -s http://localhost:8000/ | head -5  # Should return HTML
kill %1
```

- [ ] **Step 4: Commit**

```bash
git add cmd/azlake Makefile
git commit -m "feat: embed WebUI in server binary"
```

---

### Task 24: Final end-to-end test with WebUI

- [ ] **Step 1: Build everything**

```bash
make clean
make build
```

- [ ] **Step 2: Start server and test WebUI**

```bash
rm -rf /tmp/azlake
mkdir -p /tmp/azlake/data /tmp/azlake/metadata
./azlake setup --config azlake.local.yaml --user-name admin
./azlake run --config azlake.local.yaml &
sleep 3
```

Open browser to `http://localhost:8000`:
- Verify repository list page loads
- Create a repository
- Upload a file
- View file with renderer
- Create a branch
- Make changes and commit
- View diff between branches
- Test DuckDB SQL query on a parquet file (upload a test parquet first)

- [ ] **Step 3: Tag final milestone**

```bash
kill %1
git tag v0.1.0
```

---

## Phase 6 (Future): Azure Integration Testing

This phase is not part of the initial implementation but documents how to test with real Azure resources.

### Task 25: Azure integration test setup

- [ ] **Step 1: Create Azure resources**

```bash
# Create resource group
az group create --name azlake-test-rg --location eastus

# Create storage account
az storage account create \
  --name azlaketest$(date +%s | tail -c 6) \
  --resource-group azlake-test-rg \
  --location eastus \
  --sku Standard_LRS

# Create CosmosDB account (NoSQL API)
az cosmosdb create \
  --name azlaketest-cosmos \
  --resource-group azlake-test-rg \
  --kind GlobalDocumentDB \
  --default-consistency-level BoundedStaleness
```

- [ ] **Step 2: Create Azure config**

```yaml
# azlake.azure.yaml
listen_address: "0.0.0.0:8000"

auth:
  encrypt:
    secret_key: "${AZLAKE_AUTH_SECRET}"

blockstore:
  type: azure
  azure:
    storage_account: <from step 1>
    storage_access_key: <from az storage account keys list>

database:
  type: cosmosdb
  cosmosdb:
    endpoint: <from az cosmosdb show>
    database: azlake
    container: metadata
    key: <from az cosmosdb keys list>
```

- [ ] **Step 3: Run the same end-to-end test from Task 17 with Azure config**

- [ ] **Step 4: Clean up**

```bash
az group delete --name azlake-test-rg --yes
```
