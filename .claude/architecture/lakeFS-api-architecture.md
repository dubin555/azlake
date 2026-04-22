# lakeFS API Controller Architecture: Repository CRUD Operations

## Executive Summary

The lakeFS API controller implements repository CRUD operations through a multi-layered architecture:
**Controller → Catalog → Store (Graveler) → KV Store**

The dependency injection is explicit and transparent, following a clean separation of concerns where each layer is an interface that can be mocked/tested independently.

---

## 1. Controller: API Handler Layer

### File Path
`pkg/api/controller.go`

### Controller Struct Definition (lines 114-130)
```go
type Controller struct {
    Config          config.Config
    Catalog         *catalog.Catalog        // ← Primary dependency for CRUD
    Authenticator   auth.Authenticator
    Auth            auth.Service
    Authentication  authentication.Service
    BlockAdapter    block.Adapter           // ← Block storage adapter (S3, GCS, etc.)
    MetadataManager auth.MetadataManager
    Migrator        Migrator
    Collector       stats.Collector
    Actions         actionsHandler
    AuditChecker    AuditChecker
    Logger          logging.Logger
    sessionStore    sessions.Store
    PathProvider    upload.PathProvider
    usageReporter   stats.UsageReporterOperations
}
```

### CreateRepository Method (lines 2262-2343)
**Signature:**
```go
func (c *Controller) CreateRepository(
    w http.ResponseWriter,
    r *http.Request,
    body apigen.CreateRepositoryJSONRequestBody,
    params apigen.CreateRepositoryParams)
```

**Flow:**
1. **Authorization checks** via `c.authorize()`
2. **Validation** via `c.Catalog.GetRepository()` (verify not already exists)
3. **Storage validation** via `c.validateStorageNamespace()`
4. **Creates repository** by calling:
   - `c.Catalog.CreateBareRepository()` (if bare flag set) — line 2331
   - `c.Catalog.CreateRepository()` (standard) — implicit in flow
5. **Returns** HTTP 201 with Repository object

**Key call:** `c.Catalog.CreateRepository(...)`

### ListRepositories Method (lines 2198-2260)
**Signature:**
```go
func (c *Controller) ListRepositories(
    w http.ResponseWriter,
    r *http.Request,
    params apigen.ListRepositoriesParams)
```

**Flow:**
1. **Authorization** via `c.Auth.ListEffectivePolicies()` with permission filtering
2. **Fetches** repositories via `c.Catalog.ListRepositories(ctx, ...)`
3. **Transforms** catalog.Repository → apigen.Repository
4. **Returns** HTTP 200 with RepositoryList

**Key call:** `c.Catalog.ListRepositories(...opts...)`

---

## 2. Catalog: Business Logic Layer

### File Path
`pkg/catalog/catalog.go`

### Catalog Struct (lines 240-257)
```go
type Catalog struct {
    BlockAdapter            block.Adapter
    Store                   Store              // ← Core graveler interface
    managers                []io.Closer
    workPool                pond.Pool
    PathProvider            *upload.PathPartitionProvider
    BackgroundLimiter       ratelimit.Limiter
    KVStore                 kv.Store           // ← Raw KV for advanced operations
    KVStoreLimited          kv.Store           // ← Rate-limited KV store
    addressProvider         *ident.HexAddressProvider
    deleteSensor            *graveler.DeleteSensor
    UGCPrepareMaxFileSize   int64
    UGCPrepareInterval      time.Duration
    signingKey              config.SecureString
    errorToStatusCodeAndMsg ErrorToStatusCodeAndMsg
    instanceID              string
    activeTasks             stdatomic.Int64
}
```

### Store Interface (lines 143-150)
The `Catalog.Store` is an embedded Graveler that implements:
```go
type Store interface {
    graveler.KeyValueStore        // ← Entry read/write operations
    graveler.VersionController    // ← Repository lifecycle (CREATE, LIST, DELETE)
    graveler.Dumper
    graveler.Loader
    graveler.Plumbing
    graveler.Collaborator
}
```

### Repository CRUD Methods

#### CreateRepository (lines 530-555)
```go
func (c *Catalog) CreateRepository(
    ctx context.Context,
    repository string,
    storageID string,
    storageNamespace string,
    branch string,
    readOnly bool) (*Repository, error)
```
- Validates inputs (repository ID, storage namespace)
- **Delegates to:** `c.Store.CreateRepository(...)` — creates repository in KV store
- **Returns:** `*Repository` (catalog model)

#### CreateBareRepository (lines 559-590)
```go
func (c *Catalog) CreateBareRepository(
    ctx context.Context,
    repository string,
    storageID string,
    storageNamespace string,
    defaultBranchID string,
    readOnly bool) (*Repository, error)
```
- Similar validation
- **Delegates to:** `c.Store.CreateBareRepository(...)` 
- No initial branch/commit (plumbing operation)

#### GetRepository (lines 592-625)
```go
func (c *Catalog) GetRepository(ctx context.Context, repository string) (*Repository, error)
```
- **Delegates to:** `c.Store.GetRepository(...)`

#### ListRepositories (lines 676-730)
```go
func (c *Catalog) ListRepositories(
    ctx context.Context,
    limit int,
    prefix, searchString, after string,
    opts ...ListRepositoriesOptionsFunc) ([]*Repository, bool, error)
```
- Accepts options (e.g., `WithListReposPermissionFilter` from RBAC)
- **Delegates to:** `c.Store.ListRepositories(...)` for iteration
- **Returns:** `[]*Repository, hasMore, error`

---

## 3. Store/Graveler: Persistent Layer

### File Path
`pkg/graveler/graveler.go`

### VersionController Interface (lines 656-730+)
The core repository CRUD interface:

```go
type VersionController interface {
    // Create/Read
    GetRepository(ctx context.Context, repositoryID RepositoryID) 
        (*RepositoryRecord, error)
    
    CreateRepository(ctx context.Context, repositoryID RepositoryID, 
        storageID StorageID, storageNamespace StorageNamespace, 
        branchID BranchID, readOnly bool) (*RepositoryRecord, error)
    
    CreateBareRepository(ctx context.Context, repositoryID RepositoryID, 
        storageID StorageID, storageNamespace StorageNamespace, 
        defaultBranchID BranchID, readOnly bool) (*RepositoryRecord, error)
    
    ListRepositories(ctx context.Context) (RepositoryIterator, error)
    
    // Delete
    DeleteRepository(ctx context.Context, repositoryID RepositoryID, 
        opts ...SetOptionsFunc) error
    
    // Metadata
    GetRepositoryMetadata(ctx context.Context, repositoryID RepositoryID) 
        (RepositoryMetadata, error)
    
    SetRepositoryMetadata(ctx context.Context, repository *RepositoryRecord, 
        updateFunc RepoMetadataUpdateFunc) error
    
    // Branch operations...
    // Tag operations...
    // Commit operations...
}
```

### KeyValueStore Interface (lines 628-654)
For entry-level operations:
```go
type KeyValueStore interface {
    Get(ctx context.Context, repository *RepositoryRecord, ref Ref, key Key, 
        opts ...GetOptionsFunc) (*Value, error)
    GetByCommitID(ctx context.Context, repository *RepositoryRecord, 
        commitID CommitID, key Key) (*Value, error)
    GetRangeIDByKey(ctx context.Context, repository *RepositoryRecord, 
        commitID CommitID, key Key) (RangeID, error)
    Set(ctx context.Context, repository *RepositoryRecord, branchID BranchID, 
        key Key, value Value, opts ...SetOptionsFunc) error
    Update(ctx context.Context, repository *RepositoryRecord, branchID BranchID, 
        key Key, update ValueUpdateFunc, opts ...SetOptionsFunc) error
    Delete(ctx context.Context, repository *RepositoryRecord, branchID BranchID, 
        key Key, opts ...SetOptionsFunc) error
    DeleteBatch(ctx context.Context, repository *RepositoryRecord, 
        branchID BranchID, keys []Key, opts ...SetOptionsFunc) error
    List(ctx context.Context, repository *RepositoryRecord, ref Ref, 
        batchSize int) (ValueIterator, error)
}
```

### Graveler Instantiation (lines 415-423)
In `pkg/catalog/catalog.go::New()`:
```go
gStore := graveler.NewGraveler(graveler.GravelerConfig{
    CommittedManager:         committedManager,     // Persistent data layer
    StagingManager:           stagingManager,       // Staging area (KV-based)
    RefManager:               refManager,           // Refs/branches (KV-based)
    GarbageCollectionManager: gcManager,
    ProtectedBranchesManager: protectedBranchesManager,
    DeleteSensor:             deleteSensor,
    WorkPool:                 workPool,
})

cat := &Catalog{
    Store: gStore,  // ← Embedded Graveler
    KVStore: cfg.KVStore,  // ← Raw KV access
    KVStoreLimited: storeLimiter,
    ...
}
```

---

## 4. KV Store: Foundation Layer

### File Path
`pkg/kv/store.go`

### Store Interface (lines 88-111)
```go
type Store interface {
    // Get returns value or ErrNotFound
    Get(ctx context.Context, partitionKey, key []byte) (*ValueWithPredicate, error)
    
    // Set stores value (overwrite if exists)
    Set(ctx context.Context, partitionKey, key, value []byte) error
    
    // SetIf conditional set (compare-and-swap)
    SetIf(ctx context.Context, partitionKey, key, value []byte, 
        valuePredicate Predicate) error
    
    // Delete (no error if key missing)
    Delete(ctx context.Context, partitionKey, key []byte) error
    
    // Scan iterates over partition entries
    Scan(ctx context.Context, partitionKey []byte, 
        options ScanOptions) (EntriesIterator, error)
    
    // Close releases resources
    Close()
}
```

### Driver Pattern (lines 58-63)
```go
type Driver interface {
    Open(ctx context.Context, params kvparams.Config) (Store, error)
}
```

**Registered Drivers** (from `cmd/lakefs/cmd/run.go` imports):
- `postgres` — Primary production backend
- `dynamodb` — AWS DynamoDB
- `cosmosdb` — Azure Cosmos DB
- `local` — File-based (testing)
- `mem` — In-memory (testing)

### KV Initialization (from `cmd/lakefs/cmd/run.go` lines 85-93)
```go
kvParams, err := kvparams.NewConfig(&baseCfg.Database)
if err != nil {
    logger.WithError(err).Fatal("Get KV params")
}
kvStore, err := kv.Open(ctx, enableKVParamsMetrics(kvParams))
if err != nil {
    logger.WithError(err).Fatal("Failed to open KV store")
}
defer kvStore.Close()
```

---

## 5. Dependency Injection: Server Startup Path

### File Path
`cmd/lakefs/cmd/run.go`

### Initialization Chain (lines 85-234)

#### Step 1: KV Store (lines 85-93)
```go
kvStore, err := kv.Open(ctx, enableKVParamsMetrics(kvParams))
```

#### Step 2: Auth Services (lines 107-112)
```go
authService := auth.NewAuthService(ctx, cfg, logger, kvStore, authMetadataManager)
authenticationService, err := authentication.NewAuthenticationService(ctx, cfg, logger)
```

#### Step 3: Block Adapter (lines 119-127)
```go
blockStore, err := blockfactory.BuildBlockAdapterWithMetrics(ctx, bufferedCollector, cfg)
```

#### Step 4: Catalog (lines 138-149)
```go
catalogConfig := catalog.Config{
    Config:                  cfg,
    KVStore:                 kvStore,              // ← Injected
    PathProvider:            upload.DefaultPathProvider,
    ErrorToStatusCodeAndMsg: api.ErrorToStatusAndMsg,
}

c, err := catalog.New(ctx, catalogConfig)
// Internally creates Graveler with KVStore
```

#### Step 5: API Controller (lines 217-234, pkg/api/serve.go lines 85-101)
```go
controller := NewController(
    cfg,
    catalog,                 // ← Injected Catalog (which has Store)
    authenticator,
    authService,
    authenticationService,
    blockAdapter,
    metadataManager,
    migrator,
    collector,
    actions,
    auditChecker,
    logger,
    sessionStore,
    pathProvider,
    usageReporter,
)
```

---

## 6. Dependency Chain for CreateRepository

### Minimum Chain to Make CreateRepository Work

```
HTTP Request (POST /repositories)
    ↓
Controller.CreateRepository()
    ├─ checks: c.authorize()
    ├─ calls: c.Catalog.GetRepository()     [check exists]
    └─ calls: c.Catalog.CreateRepository()
        ↓
    Catalog.CreateRepository()
        ├─ validates input
        └─ calls: c.Store.CreateRepository()
            ↓
        Graveler.CreateRepository()
            ├─ creates repository record
            ├─ creates default branch
            └─ calls: KVStore operations via RefManager
                ↓
            RefManager
                ├─ stores repository metadata in KVStore
                ├─ partition key: "repository:" + repoName
                ├─ stores branch references
                └─ returns RepositoryRecord
        ↓
    Catalog returns Repository (public model)
    ↓
Controller returns HTTP 201 + Repository JSON
```

### Required Dependencies (Minimal Set)
| Component | Type | Created | Used For |
|-----------|------|---------|----------|
| `kv.Store` | Interface | `kv.Open()` | Persistent storage of all metadata |
| `graveler.Graveler` | Concrete | `graveler.NewGraveler()` | Repository lifecycle management |
| `catalog.Catalog` | Concrete | `catalog.New()` | Business logic orchestration |
| `api.Controller` | Concrete | `api.NewController()` | HTTP request handling |
| `config.Config` | Concrete | `LoadConfig()` | Configuration settings |

### Data Flow (Logical)
```
Request Body: { "name": "myrepo", "storageNamespace": "s3://bucket/repo" }
    ↓
Controller validates permissions
    ↓
Controller → Catalog → Graveler → KV Store
    ↓ (via RefManager)
KV Store writes:
    - repo:{repoName} = {RepositoryRecord}
    - repo:{repoName}:branch:main = {BranchPointer}
    - (additional metadata entries)
    ↓
Response: HTTP 201 { "id": "myrepo", "storageNamespace": "s3://...", ... }
```

---

## 7. Key Design Patterns

### 1. Explicit Interface-Based Composition
- **No global state:** All dependencies injected via constructors
- **Testable:** Each layer can be mocked (e.g., mock Store for Catalog tests)
- **Layered:** Clear separation between HTTP, business logic, persistence

### 2. Partition-Based KV Storage
- Repository metadata stored with partition key: `"repository:" + repoName`
- Enables scalability (sharding across partitions)
- Supports conditional writes (compare-and-swap)

### 3. Options Pattern (Functional Options)
```go
ListRepositories(ctx, limit, prefix, search, after, opts...)
```
- `opts` = `ListRepositoriesOptionsFunc` — allows extensibility
- Example: `WithListReposPermissionFilter(user, policies)`
- Used for RBAC filtering without API signature changes

### 4. Rate Limiting
- `KVStoreLimited` wraps raw `KVStore` with rate limiter
- `stagingManager`, `refManager` use limited store for background operations
- Critical operations use direct `KVStore` access

### 5. Error Conversion
- `ErrorToStatusCodeAndMsg` function injected into Catalog
- Maps domain errors to HTTP status codes
- Centralized error handling (pkg/api/errors.go)

---

## 8. File Structure Summary

| Path | Purpose |
|------|---------|
| `pkg/api/controller.go` | HTTP handlers + auth checks |
| `pkg/api/serve.go` | Dependency injection, router setup |
| `pkg/catalog/catalog.go` | Repository CRUD business logic |
| `pkg/graveler/graveler.go` | Persistent store interfaces |
| `pkg/graveler/ref/` | Branch/ref manager (KV operations) |
| `pkg/graveler/staging/` | Staging area manager (KV operations) |
| `pkg/kv/store.go` | KV interface + driver registry |
| `pkg/kv/postgres/` | PostgreSQL driver |
| `pkg/kv/dynamodb/` | DynamoDB driver |
| `pkg/kv/local/` | File-based driver (tests) |
| `cmd/lakefs/cmd/run.go` | Server startup + wiring |

---

## 9. Example: How CreateRepository Reaches the KV Store

```
1. HTTP POST /repositories
   ↓
2. apigen handler routes to Controller.CreateRepository()
   ↓
3. Controller checks auth, calls:
   c.Catalog.CreateRepository("myrepo", "s3", "s3://bucket/repo", "main", false)
   ↓
4. Catalog.CreateRepository validates, calls:
   c.Store.CreateRepository(repositoryID, storageID, storageNS, branchID, readOnly)
   ↓
5. Graveler.CreateRepository (Store implementation) calls RefManager:
   refManager.SetRepository(ctx, repositoryID, storageID, storageNS, branchID)
   ↓
6. RefManager calls KVStore.Set():
   kvStore.Set(ctx, partitionKey="kv:ref", key="repo:myrepo", value={...})
   kvStore.Set(ctx, partitionKey="kv:ref", key="repo:myrepo:branch:main", value={...})
   ↓
7. KV Driver (e.g., PostgreSQL) executes SQL INSERT into metadata tables
   ↓
8. Graveler returns RepositoryRecord
   ↓
9. Catalog wraps in Repository model
   ↓
10. Controller returns HTTP 201 with JSON response
```

---

## Summary Table

| Layer | File | Struct/Interface | Method | Calls |
|-------|------|------------------|--------|-------|
| **HTTP** | `pkg/api/controller.go` | `Controller` | `CreateRepository()` | `Catalog.CreateRepository()` |
| **Business Logic** | `pkg/catalog/catalog.go` | `Catalog` | `CreateRepository()` | `Store.CreateRepository()` |
| **Persistence** | `pkg/graveler/graveler.go` | `VersionController` (interface) | `CreateRepository()` | `RefManager.SetRepository()` |
| **KV Access** | `pkg/graveler/ref/ref_manager.go` | `RefManager` | `SetRepository()` | `KVStore.Set()` |
| **Storage** | `pkg/kv/store.go` | `Store` (interface) | `Set()` | Driver-specific (SQL, DynamoDB, etc.) |

