package kvparams

import (
	"net/http"
	"time"
)

type Config struct {
	Type     string
	Postgres *Postgres
	DynamoDB *DynamoDB
	Local    *Local
	CosmosDB *CosmosDB
	Redis    *Redis
}

type Local struct {
	Path          string
	SyncWrites    bool
	PrefetchSize  int
	EnableLogging bool
}

type Postgres struct {
	ConnectionString      string
	MaxOpenConnections    int32
	MaxIdleConnections    int32
	ConnectionMaxLifetime time.Duration
	ScanPageSize          int
	Metrics               bool
}

type DynamoDB struct {
	TableName                                  string
	ScanLimit                                  int64
	MaxAttempts                                int
	Endpoint                                   string
	AwsRegion                                  string
	AwsProfile                                 string
	AwsAccessKeyID                             string
	AwsSecretAccessKey                         string
	HealthCheckInterval                        time.Duration
	MaxConnectionsPerHost                      int
	CredentialsCacheExpiryWindow               time.Duration
	CredentialsCacheExpiryWindowJitterFraction float64
}

type CosmosDB struct {
	Key               string
	Endpoint          string
	Database          string
	Container         string
	Throughput        int32
	Autoscale         bool
	Client            *http.Client
	StrongConsistency bool
}

type Redis struct {
	Endpoint      string
	Username      string
	Password      string
	Database      int
	PoolSize      int
	MinIdleConns  int
	DialTimeout   time.Duration
	ReadTimeout   time.Duration
	WriteTimeout  time.Duration
	Namespace     string
	EnableTLS     bool
	ClusterMode   bool
	TLSSkipVerify bool
	BatchSize     int
}
