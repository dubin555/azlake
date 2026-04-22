# azlake

**Azure-native data version control. Git for your data lake.**

azlake brings Git-like version control to your data stored on Azure. Create repositories, branches, commits, and diffs — all backed by **Azure Blob Storage** and **Azure CosmosDB**.

Extracted and adapted from [lakeFS](https://github.com/treeverse/lakeFS), rebuilt for Azure-first workflows.

## Features

- **Git-like operations** — repositories, branches, commits, diffs, merges
- **Azure Blob Storage** — objects stored directly in Azure Blob with SAS URL support
- **Azure CosmosDB** — metadata stored in CosmosDB (serverless-friendly)
- **Local mode** — run entirely locally with BadgerDB + filesystem (zero cloud dependency)
- **Web UI** — browse repos, view files, upload objects, inspect commits
- **DuckDB WASM** — query Parquet/CSV files directly in the browser via SQL
- **REST API** — OpenAPI-compatible API for all operations

## Architecture

```
┌─────────────┐     ┌──────────────────────────────────┐
│   Web UI    │────▶│         REST API (Go)             │
│  (React)    │     │                                   │
└─────────────┘     └──────┬───────────────┬────────────┘
                           │               │
                    ┌──────▼──────┐ ┌──────▼──────────┐
                    │  Metadata   │ │  Object Storage  │
                    │             │ │                   │
                    │ • BadgerDB  │ │ • Local FS        │
                    │ • CosmosDB  │ │ • Azure Blob      │
                    └─────────────┘ └──────────────────┘
```

## Quick Start

### Local mode (no Azure required)

```bash
go build -o azlake ./cmd/azlake/
./azlake run
```

Open http://localhost:8000 — data is stored in `~/.azlake/`.

### Azure mode

```bash
export AZLAKE_STORAGE_BACKEND=azure
export AZURE_STORAGE_ACCOUNT=<your-account>
export AZURE_STORAGE_KEY=<your-key>          # optional: falls back to az login
export AZURE_STORAGE_CONTAINER=<container>

export AZLAKE_KV_BACKEND=cosmosdb
export COSMOS_ENDPOINT=<your-endpoint>
export COSMOS_KEY=<your-key>                 # optional: falls back to az login
export COSMOS_DATABASE=<database>
export COSMOS_CONTAINER=<container>

./azlake run
```

## Configuration

| Environment Variable | Default | Description |
|---|---|---|
| `AZLAKE_STORAGE_BACKEND` | `local` | Object storage: `local` or `azure` |
| `AZLAKE_KV_BACKEND` | `badger` | Metadata store: `badger` or `cosmosdb` |
| `AZURE_STORAGE_ACCOUNT` | — | Azure Storage account name |
| `AZURE_STORAGE_KEY` | — | Azure Storage key (optional, uses DefaultAzureCredential) |
| `AZURE_STORAGE_CONTAINER` | — | Azure Blob container name |
| `COSMOS_ENDPOINT` | — | CosmosDB endpoint URL |
| `COSMOS_KEY` | — | CosmosDB key (optional, uses DefaultAzureCredential) |
| `COSMOS_DATABASE` | — | CosmosDB database name |
| `COSMOS_CONTAINER` | — | CosmosDB container name |

## API Examples

```bash
# Create a repository
curl -X POST http://localhost:8000/api/v1/repositories \
  -H 'Content-Type: application/json' \
  -d '{"name":"my-repo","storage_namespace":"az://container/my-repo","default_branch":"main"}'

# Upload a file
curl -X POST "http://localhost:8000/api/v1/repositories/my-repo/branches/main/objects?path=data.csv" \
  -H 'Content-Type: application/octet-stream' \
  --data-binary @data.csv

# Commit
curl -X POST "http://localhost:8000/api/v1/repositories/my-repo/branches/main/commits" \
  -H 'Content-Type: application/json' \
  -d '{"message":"Add data.csv"}'

# Read a file
curl "http://localhost:8000/api/v1/repositories/my-repo/refs/main/objects?path=data.csv"
```

## Building

### Prerequisites

- Go 1.21+
- Node.js 18+ (for Web UI)

### Build

```bash
# Build frontend
cd webui && npm install && npm run build && cd ..

# Build backend (embeds frontend)
go build -o azlake ./cmd/azlake/
```

## Project Structure

```
cmd/azlake/          — CLI entrypoint
pkg/azcat/           — Core catalog (repos, branches, commits, objects)
  catalog.go         — Catalog operations
  storage.go         — ObjectStorage interface + LocalStorage
  storage_azure.go   — Azure Blob Storage implementation
  kvstore.go         — BadgerDB KV store
  kvstore_cosmos.go  — CosmosDB KV store
pkg/api/             — REST API handlers
webui/               — React frontend
```

## License

Apache 2.0 — see [LICENSE](LICENSE).

## Acknowledgments

Built on the foundation of [lakeFS](https://github.com/treeverse/lakeFS) by Treeverse.
