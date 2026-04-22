# azlake

**Azure 原生的数据版本控制系统。给你的数据湖加上 Git。**

azlake 为存储在 Azure 上的数据提供类 Git 的版本控制能力。支持创建仓库、分支、提交、差异对比 —— 全部由 **Azure Blob Storage** 和 **Azure CosmosDB** 驱动。

基于 [lakeFS](https://github.com/treeverse/lakeFS) 提取和改造，专为 Azure 生态系统重构。

## 功能特性

- **类 Git 操作** — 仓库、分支、提交、差异对比、合并
- **Azure Blob Storage** — 对象直接存储在 Azure Blob 中，支持 SAS URL 直链访问
- **Azure CosmosDB** — 元数据存储在 CosmosDB 中（支持 Serverless 模式）
- **本地模式** — 使用 BadgerDB + 本地文件系统运行，零云依赖
- **Web 界面** — 浏览仓库、查看文件、上传对象、查看提交历史
- **DuckDB WASM** — 在浏览器中用 SQL 直接查询 Parquet/CSV 文件
- **REST API** — 兼容 OpenAPI 规范的完整 API

## 架构

```
┌─────────────┐     ┌──────────────────────────────────┐
│   Web UI    │────▶│         REST API (Go)             │
│  (React)    │     │                                   │
└─────────────┘     └──────┬───────────────┬────────────┘
                           │               │
                    ┌──────▼──────┐ ┌──────▼──────────┐
                    │   元数据    │ │    对象存储       │
                    │             │ │                   │
                    │ • BadgerDB  │ │ • 本地文件系统     │
                    │ • CosmosDB  │ │ • Azure Blob      │
                    └─────────────┘ └──────────────────┘
```

## 快速开始

### 本地模式（无需 Azure）

```bash
go build -o azlake ./cmd/azlake/
./azlake run
```

打开 http://localhost:8000 — 数据存储在 `~/.azlake/` 目录。

### Azure 模式

```bash
export AZLAKE_STORAGE_BACKEND=azure
export AZURE_STORAGE_ACCOUNT=<你的存储账户>
export AZURE_STORAGE_KEY=<你的密钥>          # 可选：不填则使用 az login 凭证
export AZURE_STORAGE_CONTAINER=<容器名>

export AZLAKE_KV_BACKEND=cosmosdb
export COSMOS_ENDPOINT=<你的端点>
export COSMOS_KEY=<你的密钥>                 # 可选：不填则使用 az login 凭证
export COSMOS_DATABASE=<数据库名>
export COSMOS_CONTAINER=<容器名>

./azlake run
```

## 配置项

| 环境变量 | 默认值 | 说明 |
|---|---|---|
| `AZLAKE_STORAGE_BACKEND` | `local` | 对象存储后端：`local` 或 `azure` |
| `AZLAKE_KV_BACKEND` | `badger` | 元数据存储后端：`badger` 或 `cosmosdb` |
| `AZURE_STORAGE_ACCOUNT` | — | Azure 存储账户名 |
| `AZURE_STORAGE_KEY` | — | Azure 存储密钥（可选，支持 DefaultAzureCredential） |
| `AZURE_STORAGE_CONTAINER` | — | Azure Blob 容器名 |
| `COSMOS_ENDPOINT` | — | CosmosDB 端点 URL |
| `COSMOS_KEY` | — | CosmosDB 密钥（可选，支持 DefaultAzureCredential） |
| `COSMOS_DATABASE` | — | CosmosDB 数据库名 |
| `COSMOS_CONTAINER` | — | CosmosDB 容器名 |

## API 示例

```bash
# 创建仓库
curl -X POST http://localhost:8000/api/v1/repositories \
  -H 'Content-Type: application/json' \
  -d '{"name":"my-repo","storage_namespace":"az://container/my-repo","default_branch":"main"}'

# 上传文件
curl -X POST "http://localhost:8000/api/v1/repositories/my-repo/branches/main/objects?path=data.csv" \
  -H 'Content-Type: application/octet-stream' \
  --data-binary @data.csv

# 提交
curl -X POST "http://localhost:8000/api/v1/repositories/my-repo/branches/main/commits" \
  -H 'Content-Type: application/json' \
  -d '{"message":"添加 data.csv"}'

# 读取文件
curl "http://localhost:8000/api/v1/repositories/my-repo/refs/main/objects?path=data.csv"
```

## 构建

### 前置条件

- Go 1.21+
- Node.js 18+（Web UI）

### 编译

```bash
# 构建前端
cd webui && npm install && npm run build && cd ..

# 构建后端（内嵌前端资源）
go build -o azlake ./cmd/azlake/
```

## 项目结构

```
cmd/azlake/          — CLI 入口
pkg/azcat/           — 核心 Catalog（仓库、分支、提交、对象）
  catalog.go         — Catalog 操作
  storage.go         — ObjectStorage 接口 + 本地存储实现
  storage_azure.go   — Azure Blob Storage 实现
  kvstore.go         — BadgerDB KV 存储
  kvstore_cosmos.go  — CosmosDB KV 存储
pkg/api/             — REST API 处理器
webui/               — React 前端
```

## 许可证

Apache 2.0 — 详见 [LICENSE](LICENSE)。

## 致谢

基于 [lakeFS](https://github.com/treeverse/lakeFS)（Treeverse）构建。
