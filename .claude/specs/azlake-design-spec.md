# azlake Design Spec

**Project**: azlake — Azure-native data version control system
**Module**: `github.com/dubin555/azlake`
**Approach**: Core extraction from lakeFS (new project, copy needed packages)
**Date**: 2026-04-21

## 1. Overview

azlake is a stripped-down fork of [lakeFS](https://github.com/treeverse/lakefs) that removes multi-cloud support and unnecessary features, keeping only Azure-native functionality. The goal is an open-source, lightweight data version control system purpose-built for the Azure ecosystem.

### What it does

- Git-like version control for data stored in Azure Blob Storage (branch, commit, merge, revert, tag)
- Metadata stored in Azure CosmosDB (production) or local BadgerDB (development)
- REST API + CLI for all operations
- Web UI with DuckDB WASM SQL query explorer and file preview
- Webhook hooks for pre/post-commit and pre/post-merge events

### What it removes from lakeFS

| Removed | Lines saved | Reason |
|---------|-------------|--------|
| S3/GCS/Transient/Mem block adapters | ~3,500 | Azure-only |
| DynamoDB/PostgreSQL/Mem KV drivers | ~1,400 | CosmosDB + Local only |
| S3 Gateway | ~6,300 | No S3 compatibility needed |
| Lua VM + Airflow hooks | ~4,300 | Only simple webhook needed |
| RBAC/Policy engine | ~9,000 | API Key auth only |
| Pull Request system | ~1,000 | Not needed |
| Full React WebUI | ~19,300 | Replaced by slim UI |
| All SDKs (Python/Java/Rust/Spark) | ~160,000 | Generated, not needed |
| Hadoop FileSystem plugin | ~4,000 | No Spark integration |
| **Total removed** | **~209,000** | |

### Target code size

| Component | Estimated lines |
|-----------|----------------|
| Graveler (version engine) | ~13,500 |
| KV core + CosmosDB + Local | ~2,700 |
| Block core + Azure + Local | ~3,400 |
| Catalog | ~6,300 |
| REST API (trimmed endpoints) | ~6,000 |
| Auth (API Key only) | ~2,500 |
| Webhook hooks | ~1,000 |
| Config | ~1,000 |
| CLI (lakectl equivalent) | ~5,000 |
| Supporting packages | ~2,000 |
| **Go backend total** | **~43,400** |
| Web UI (slim React) | ~5,000 |
| **Grand total** | **~51,000** |

## 2. Architecture

```
┌─────────────────────────────────────────────┐
│           WebUI (slim React + Vite)          │
│   DuckDB WASM SQL | Diff View | File Preview│
├─────────────────────────────────────────────┤
│         REST API (go-chi/chi router)         │
│         + Pre-signed URL endpoints           │
├─────────────────────────────────────────────┤
│            Catalog (coordination)            │
├──────────┬───────────┬──────────────────────┤
│ Graveler │   Auth    │  Webhook Hooks       │
│ (version │ (API Key  │ (HTTP callback only, │
│  engine) │  only)    │  pre/post commit +   │
│          │           │  pre/post merge)     │
├──────────┴───────────┴──────────────────────┤
│   KV Store interface                         │
│     CosmosDB (production)                    │
│     BadgerDB (local dev/test)                │
├─────────────────────────────────────────────┤
│   Block Adapter interface                    │
│     Azure Blob Storage (production)          │
│     Local filesystem (local dev/test)        │
└─────────────────────────────────────────────┘
```

### Data flow: file read

```
Client request: GET /api/v1/repositories/repo/refs/main/objects?path=data/file.parquet

  1. REST API layer receives request
  2. Auth: validate API key
  3. Catalog.GetEntry() → Graveler.Get()
     a. Check staging area (KV: current token → sealed tokens)
     b. Check committed MetaRange (Block: read SSTable index)
     c. Return: { PhysicalAddress, Size, Checksum, ContentType }
  4. Two modes:
     a. Proxy: BlockStore.Get(physicalAddress) → stream to client
     b. Pre-signed: BlockStore.GetPreSignedURL(physicalAddress) → return SAS URL
```

### Data flow: file write + commit

```
1. Upload: PUT /api/v1/repositories/repo/branches/main/objects?path=data/file.parquet
   → BlockStore.Put() writes to Azure Blob with unique physical address
   → StagingManager.Set() writes path→physicalAddress mapping to KV

2. Commit: POST /api/v1/repositories/repo/branches/main/commits
   → Seal current staging token (no more writes to it)
   → Build new MetaRange by merging staging changes into parent MetaRange
   → Create CommitData record in KV
   → Update BranchData.commit_id in KV (CAS via ETag)
   → Fire post-commit webhook (if configured)
```

### Storage model

All "files" in azlake are virtual — the logical path and physical blob are decoupled:

```
Logical view:                          Physical reality (Azure Blob Storage):
  repo/main/users.parquet              container/a1b2c3d4e5_xid  (version 1)
  repo/main/users.parquet (updated)    container/f6g7h8i9j0_xid  (version 2)
                                       ↑ Both blobs exist until GC runs

MetaRange (commit 1): users.parquet → a1b2c3d4e5_xid
MetaRange (commit 2): users.parquet → f6g7h8i9j0_xid
```

Each file update creates a new physical blob. Old blobs are retained for version history and cleaned up by garbage collection based on retention rules.

## 3. Backend Modules

### 3.1 Graveler (version engine) — copy from lakeFS

The core version engine. Copy `pkg/graveler/` with these removals:

- Remove PullRequest-related code (`PullRequestData` proto, PR manager, PR iterators)
- Remove branch protection manager (not needed without RBAC)
- Keep: Branch, Commit, Merge, Revert, Tag, Reset, Diff, Log
- Keep: CombinedIterator, FilterTombstoneIterator (read-path merge logic)
- Keep: StagingManager, CommittedManager, RefManager interfaces

### 3.2 KV Store — two implementations

**Interface** (copy `pkg/kv/store.go`):
```go
type Store interface {
    Get(ctx, partitionKey, key) → (*ValueWithPredicate, error)
    Set(ctx, partitionKey, key, value) → error
    SetIf(ctx, partitionKey, key, value, predicate) → error
    Delete(ctx, partitionKey, key) → error
    Scan(ctx, partitionKey, options) → (EntriesIterator, error)
    Close()
}
```

**CosmosDB** (copy `pkg/kv/cosmosdb/`):
- Uses ETag for CAS (optimistic locking)
- Base32 HexEncoding for order-preserving keys
- Partition key = `{repoID}-{instanceUID}` for repo data, staging token for staged entries
- ~580 lines

**Local/BadgerDB** (copy `pkg/kv/local/`):
- For local development and testing
- ~235 lines

### 3.3 Block Adapter — two implementations

**Interface** (copy `pkg/block/` core interfaces):
- Get, Put, Delete, GetPreSignedURL, GetRange, GetProperties, Walk

**Azure Blob Storage** (copy `pkg/block/azure/`):
- Client cache for multi-account support
- SAS token generation for pre-signed URLs
- Supports standard, China, and US Government Azure domains
- ~1,800 lines

**Local filesystem** (copy `pkg/block/local/`):
- For local development and testing
- ~624 lines

### 3.4 Auth — API Key only

Drastically simplified from lakeFS's full RBAC system:

- Single-user or multi-user via API key pairs (access_key_id + secret_access_key)
- Keys stored in KV store
- No groups, no policies, no ACLs
- Middleware validates API key on each request
- Estimated ~2,500 lines (down from 12,500)

### 3.5 Webhook Hooks

Extracted from lakeFS `pkg/actions/`, keeping only the webhook type:

- Hook points: pre-commit, post-commit, pre-merge, post-merge
- Hook type: HTTP POST to configured URL
- Configuration via `_azlake_actions/` directory in repository
- Remove: Lua VM, Airflow integration, action run history storage
- Estimated ~1,000 lines (down from 6,300)

### 3.6 Catalog

Copy `pkg/catalog/` as-is — it's the coordination layer between REST API and Graveler. May need minor adjustments to remove references to deleted features (pull requests, branch protection).

### 3.7 REST API

Copy `pkg/api/` with endpoint removals:
- Remove: all auth management endpoints (users, groups, policies)
- Remove: pull request endpoints
- Remove: S3 gateway endpoints
- Keep: repositories, branches, commits, tags, objects, diff, merge, revert, config
- Add: pre-signed URL endpoint for WebUI DuckDB integration

### 3.8 CLI

Create `azlake` (server) and `azlakectl` (client) binaries:
- `azlake` replaces `lakefs` server binary
- `azlakectl` replaces `lakectl` CLI tool
- Remove commands related to deleted features

### 3.9 Factory / Config

Simplified factory with two-way switch:

```go
func BuildBlockAdapter(ctx, config) (block.Adapter, error) {
    switch config.BlockstoreType() {
    case "azure":
        return azure.NewAdapter(ctx, config.AzureParams())
    case "local":
        return local.NewAdapter(config.LocalPath())
    default:
        return nil, fmt.Errorf("unsupported blockstore type: %s", typ)
    }
}

func BuildKVStore(ctx, config) (kv.Store, error) {
    switch config.DatabaseType() {
    case "cosmosdb":
        return cosmosdb.Open(ctx, config.CosmosDBParams())
    case "local":
        return local.Open(ctx, config.LocalParams())
    default:
        return nil, fmt.Errorf("unsupported database type: %s", typ)
    }
}
```

Configuration file (`azlake.yaml`):

```yaml
# Local development
listen_address: "0.0.0.0:8000"

auth:
  encrypt:
    secret_key: "a-random-secret-key"

blockstore:
  type: local
  local:
    path: /tmp/azlake/data

database:
  type: local
  local:
    path: /tmp/azlake/metadata
```

```yaml
# Production (Azure)
listen_address: "0.0.0.0:8000"

auth:
  encrypt:
    secret_key: "${AZLAKE_AUTH_SECRET}"

blockstore:
  type: azure
  azure:
    storage_account: myaccount
    storage_access_key: "${AZURE_STORAGE_KEY}"
    # Or use managed identity (no key needed)

database:
  type: cosmosdb
  cosmosdb:
    endpoint: "https://myaccount.documents.azure.com"
    database: azlake
    container: metadata
    # Key or managed identity
```

## 4. Web UI

Slim React application built with Vite, extracted from lakeFS WebUI with significant reduction.

### 4.1 Pages

| Page | Purpose | Source |
|------|---------|--------|
| Repositories list | Create/delete repos | New (simplified) |
| Object browser | Navigate file tree, upload/download | Extract from lakeFS |
| File preview | Render file contents based on type | Extract renderers |
| DuckDB SQL query | In-browser SQL on Parquet/CSV | Extract + modify |
| Commit history | View commit log for a branch | Extract from lakeFS |
| Branch management | Create/delete/compare branches | Extract from lakeFS |
| Diff view | Compare two refs, text diff | Extract from lakeFS |
| Tags | Create/delete tags | New (simplified) |

### 4.2 DuckDB WASM Integration

The original lakeFS DuckDB integration uses S3 protocol through the S3 Gateway. Since azlake has no S3 Gateway, this is modified to use pre-signed URLs:

```
Original (lakeFS):
  DuckDB WASM → S3 protocol → lakeFS S3 Gateway → Azure Blob

azlake:
  DuckDB WASM → HTTP protocol → Pre-signed Azure SAS URL → Azure Blob direct
```

Implementation change in `duckdb.tsx`:

```typescript
// Instead of S3 endpoint configuration, fetch pre-signed URL from API
async function getLakeFSPresignedUrl(repo: string, ref: string, path: string): Promise<string> {
    const response = await fetch(
        `/api/v1/repositories/${repo}/refs/${ref}/objects/presign?path=${encodeURIComponent(path)}`
    );
    const data = await response.json();
    return data.url;  // Azure SAS URL
}

// Register with HTTP protocol instead of S3
const presignedUrl = await getLakeFSPresignedUrl(repo, ref, filePath);
db.registerFileURL(fileName, presignedUrl, DuckDBDataProtocol.HTTP, true);
```

### 4.3 File Renderers

Extracted from lakeFS `webui/src/pages/repositories/repository/fileRenderers/` (~730 lines):

| Renderer | File types | Implementation |
|----------|-----------|----------------|
| DuckDB SQL | .parquet, .csv, .tsv | `duckdb.tsx` + `data.tsx` (298 lines) |
| Code/Text | .json, .yaml, .py, .go, etc. | Syntax highlighting via Prism.js |
| Markdown | .md | GFM rendering with image URI rewriting |
| Image | .png, .jpg, .gif, .webp, .bmp | Direct display |
| GeoJSON | .geojson | Leaflet map visualization |

### 4.4 Diff Visualization

Extracted from lakeFS:
- **Branch compare**: side-by-side tree diff showing added/removed/changed files
- **Text diff**: `react-diff-viewer-continued` for line-level diff of readable files (up to 120KB)
- **Image diff**: side-by-side display of two versions

### 4.5 Not included

- User/group management pages (API Key only, no UI needed)
- Pull request pages
- Actions/hooks monitoring (webhook is too simple for a UI)
- Settings page (use config file)

## 5. Testing Strategy

### 5.1 Local testing (no Azure needed)

```bash
# Start azlake in local mode
azlake run --config local.yaml

# Create repo, upload file, commit — all local
azlakectl repo create lakefs://test-repo --storage local:///tmp/azlake/data
azlakectl fs upload lakefs://test-repo/main/test.parquet --source ./test.parquet
azlakectl commit lakefs://test-repo/main -m "initial data"
azlakectl branch create lakefs://test-repo/dev --source lakefs://test-repo/main
```

### 5.2 Azure integration testing

Requires `az login` and a test storage account:

```bash
# Create test resources
az storage account create --name azlaketest --resource-group test-rg --location eastus
az cosmosdb create --name azlaketest-cosmos --resource-group test-rg

# Run integration tests
AZLAKE_TEST_STORAGE_ACCOUNT=azlaketest \
AZLAKE_TEST_COSMOS_ENDPOINT=https://azlaketest-cosmos.documents.azure.com \
go test ./... -tags=integration
```

### 5.3 Test scope

- Unit tests: copy relevant tests from lakeFS for retained modules
- Integration tests: Azure Blob + CosmosDB end-to-end flows
- WebUI: basic Playwright smoke tests for DuckDB query and diff view

## 6. Project Structure

```
github.com/dubin555/azlake/
├── cmd/
│   ├── azlake/          # Server binary
│   └── azlakectl/       # CLI binary
├── pkg/
│   ├── api/             # REST API handlers
│   ├── auth/            # API Key authentication (simplified)
│   ├── block/
│   │   ├── azure/       # Azure Blob Storage adapter
│   │   ├── local/       # Local filesystem adapter
│   │   └── factory/     # Adapter factory (azure | local)
│   ├── catalog/         # Business coordination layer
│   ├── config/          # Configuration management
│   ├── graveler/        # Version engine (core)
│   │   ├── committed/   # MetaRange/Range management
│   │   ├── ref/         # Branch/Commit/Tag managers
│   │   └── staging/     # Staging area manager
│   ├── hooks/           # Webhook-only hooks (simplified)
│   ├── kv/
│   │   ├── cosmosdb/    # CosmosDB driver
│   │   └── local/       # BadgerDB driver
│   └── upload/          # Blob upload utilities
├── webui/
│   ├── src/
│   │   ├── components/  # Shared UI components
│   │   ├── pages/
│   │   │   ├── repositories/  # Repo browser, object viewer
│   │   │   ├── compare/       # Diff visualization
│   │   │   └── commits/       # Commit history
│   │   ├── renderers/         # File preview (DuckDB, Markdown, etc.)
│   │   └── lib/               # API client, utilities
│   ├── package.json
│   └── vite.config.ts
├── api/
│   └── openapi.yaml     # API spec (trimmed)
├── go.mod
├── go.sum
├── Makefile
├── Dockerfile
└── azlake.yaml.example
```

## 7. Go Module Dependencies

### Keep

| Dependency | Purpose |
|------------|---------|
| `github.com/go-chi/chi` | HTTP router |
| `github.com/spf13/cobra` | CLI framework |
| `github.com/spf13/viper` | Configuration |
| `github.com/Azure/azure-sdk-for-go/sdk/storage/azblob` | Azure Blob |
| `github.com/Azure/azure-sdk-for-go/sdk/azidentity` | Azure auth |
| `github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos` | CosmosDB |
| `github.com/dgraph-io/badger/v4` | Local KV store |
| `github.com/cockroachdb/pebble` | MetaRange SSTable engine |
| `google.golang.org/protobuf` | Protobuf serialization |
| `github.com/rs/xid` | Unique ID generation |
| `github.com/cenkalti/backoff/v4` | Retry with exponential backoff |

### Remove

| Dependency | Reason |
|------------|--------|
| `github.com/aws/aws-sdk-go-v2` | No S3 |
| `cloud.google.com/go/storage` | No GCS |
| `github.com/Shopify/go-lua` | No Lua VM |
| `github.com/xitongsys/parquet-go` | Only needed for Spark |
| Various auth/OIDC libraries | API Key only |

## 8. Migration Path from lakeFS

For each module, the extraction process:

1. **Copy** the package directory into the new project
2. **Rename** imports from `github.com/treeverse/lakefs/` to `github.com/dubin555/azlake/`
3. **Delete** files related to removed features
4. **Simplify** factory/switch statements to only Azure + Local
5. **Compile** and fix any broken references
6. **Test** with local mode first, then Azure integration

Order of extraction (dependencies flow downward):

```
Phase 1: Foundation (no dependencies on other azlake packages)
  pkg/kv/          → KV interface + CosmosDB + Local drivers
  pkg/block/       → Block interface + Azure + Local adapters
  pkg/config/      → Configuration

Phase 2: Core engine (depends on kv + block)
  pkg/graveler/    → Version engine
  pkg/upload/      → Blob upload utilities

Phase 3: Business layer (depends on graveler)
  pkg/catalog/     → Catalog coordination
  pkg/auth/        → API Key auth (rewrite, simplified)
  pkg/hooks/       → Webhook hooks (extract from actions)

Phase 4: Interface layer (depends on catalog + auth)
  pkg/api/         → REST API handlers
  cmd/azlake/      → Server binary
  cmd/azlakectl/   → CLI binary

Phase 5: Web UI
  webui/           → Slim React app with DuckDB + renderers + diff
```

## 9. Non-Goals

Things explicitly out of scope for v1:

- S3 compatibility layer
- Multi-user RBAC / fine-grained permissions
- Spark / Hadoop integration
- Import from existing Azure storage (can be added later)
- Garbage collection automation (manual cleanup for v1)
- High availability / clustering (single instance for v1)
