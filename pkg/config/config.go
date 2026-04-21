package config

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/go-viper/mapstructure/v2"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	apiparams "github.com/dubin555/azlake/pkg/api/params"
	blockparams "github.com/dubin555/azlake/pkg/block/params"
	"github.com/dubin555/azlake/pkg/logging"
)

var (
	ErrBadConfiguration    = errors.New("bad configuration")
	ErrBadDomainNames      = fmt.Errorf("%w: domain names are prefixes", ErrBadConfiguration)
	ErrMissingRequiredKeys = fmt.Errorf("%w: missing required keys", ErrBadConfiguration)
	ErrNoStorageConfig     = errors.New("no storage config")
)

const (
	UseLocalConfiguration   = "local-settings"
	QuickstartConfiguration = "quickstart"
	SingleBlockstoreID      = ""
)

const (
	AuthRBACNone       = "none"
	AuthRBACSimplified = "simplified"
	AuthRBACExternal   = "external"
	AuthRBACInternal   = "internal"
)

type Logging struct {
	Format        string   `mapstructure:"format"`
	Level         string   `mapstructure:"level"`
	Output        []string `mapstructure:"output"`
	FileMaxSizeMB int      `mapstructure:"file_max_size_mb"`
	FilesKeep     int      `mapstructure:"files_keep"`
	AuditLogLevel string   `mapstructure:"audit_log_level"`
	TraceRequestHeaders bool `mapstructure:"trace_request_headers"`
}

// Database - holds metadata KV configuration (Azure-focused: CosmosDB + Local)
type Database struct {
	DropTables bool   `mapstructure:"drop_tables"`
	Type       string `mapstructure:"type" validate:"required"`

	Local *struct {
		Path          string `mapstructure:"path"`
		SyncWrites    bool   `mapstructure:"sync_writes"`
		PrefetchSize  int    `mapstructure:"prefetch_size"`
		EnableLogging bool   `mapstructure:"enable_logging"`
	} `mapstructure:"local"`

	CosmosDB *struct {
		Key        SecureString `mapstructure:"key"`
		Endpoint   string       `mapstructure:"endpoint"`
		Database   string       `mapstructure:"database"`
		Container  string       `mapstructure:"container"`
		Throughput int32        `mapstructure:"throughput"`
		Autoscale  bool         `mapstructure:"autoscale"`
	} `mapstructure:"cosmosdb"`
}

// ApproximatelyCorrectOwnership configures an approximate ownership.
type ApproximatelyCorrectOwnership struct {
	Enabled bool          `mapstructure:"enabled"`
	Refresh time.Duration `mapstructure:"refresh"`
	Acquire time.Duration `mapstructure:"acquire"`
}

// AdapterConfig configures a blockstore adapter.
type AdapterConfig interface {
	BlockstoreType() string
	BlockstoreDescription() string
	BlockstoreLocalParams() (blockparams.Local, error)
	BlockstoreAzureParams() (blockparams.Azure, error)
	GetDefaultNamespacePrefix() *string
	IsBackwardsCompatible() bool
	ID() string
}

type BlockstoreLocal struct {
	Path                    string   `mapstructure:"path"`
	ImportEnabled           bool     `mapstructure:"import_enabled"`
	ImportHidden            bool     `mapstructure:"import_hidden"`
	AllowedExternalPrefixes []string `mapstructure:"allowed_external_prefixes"`
}

type BlockstoreAzure struct {
	TryTimeout       time.Duration `mapstructure:"try_timeout"`
	StorageAccount   string        `mapstructure:"storage_account"`
	StorageAccessKey string        `mapstructure:"storage_access_key"`
	AuthMethodDeprecated string    `mapstructure:"auth_method"`
	PreSignedExpiry      time.Duration `mapstructure:"pre_signed_expiry"`
	DisablePreSigned     bool          `mapstructure:"disable_pre_signed"`
	DisablePreSignedUI   bool          `mapstructure:"disable_pre_signed_ui"`
	ChinaCloudDeprecated bool   `mapstructure:"china_cloud"`
	TestEndpointURL      string `mapstructure:"test_endpoint_url"`
	Domain               string `mapstructure:"domain"`
}

type Blockstore struct {
	Signing struct {
		SecretKey SecureString `mapstructure:"secret_key"`
	} `mapstructure:"signing"`
	Type                   string           `mapstructure:"type"`
	DefaultNamespacePrefix *string          `mapstructure:"default_namespace_prefix"`
	Local                  *BlockstoreLocal `mapstructure:"local"`
	Azure                  *BlockstoreAzure `mapstructure:"azure"`
}

func (b *Blockstore) GetStorageIDs() []string {
	return []string{SingleBlockstoreID}
}

func (b *Blockstore) GetStorageByID(id string) AdapterConfig {
	if id != SingleBlockstoreID {
		return nil
	}
	return b
}

func (b *Blockstore) BlockstoreType() string {
	return b.Type
}

func (b *Blockstore) BlockstoreLocalParams() (blockparams.Local, error) {
	localPath := b.Local.Path
	path, err := homedir.Expand(localPath)
	if err != nil {
		return blockparams.Local{}, fmt.Errorf("parse blockstore location URI %s: %w", localPath, err)
	}
	params := blockparams.Local(*b.Local)
	params.Path = path
	return params, nil
}

func (b *Blockstore) BlockstoreAzureParams() (blockparams.Azure, error) {
	if b.Azure.AuthMethodDeprecated != "" {
		logging.ContextUnavailable().Warn("blockstore.azure.auth_method is deprecated. Value is no longer used.")
	}
	if b.Azure.ChinaCloudDeprecated {
		logging.ContextUnavailable().Warn("blockstore.azure.china_cloud is deprecated. Please pass Domain = 'blob.core.chinacloudapi.cn'")
		b.Azure.Domain = "blob.core.chinacloudapi.cn"
	}
	return blockparams.Azure{
		StorageAccount:     b.Azure.StorageAccount,
		StorageAccessKey:   b.Azure.StorageAccessKey,
		TryTimeout:         b.Azure.TryTimeout,
		PreSignedExpiry:    b.Azure.PreSignedExpiry,
		TestEndpointURL:    b.Azure.TestEndpointURL,
		Domain:             b.Azure.Domain,
		DisablePreSigned:   b.Azure.DisablePreSigned,
		DisablePreSignedUI: b.Azure.DisablePreSignedUI,
	}, nil
}

func (b *Blockstore) BlockstoreDescription() string {
	return ""
}

func (b *Blockstore) GetDefaultNamespacePrefix() *string {
	return b.DefaultNamespacePrefix
}

func (b *Blockstore) IsBackwardsCompatible() bool {
	return false
}

func (b *Blockstore) ID() string {
	return SingleBlockstoreID
}

func (b *Blockstore) SigningKey() SecureString {
	return b.Signing.SecretKey
}

// GetActualStorageID returns the actual storageID of the storage
func GetActualStorageID(storageConfig StorageConfig, storageID string) string {
	if storageID == SingleBlockstoreID {
		if storage := storageConfig.GetStorageByID(SingleBlockstoreID); storage != nil {
			return storage.ID()
		}
	}
	return storageID
}

type Config interface {
	GetBaseConfig() *BaseConfig
	StorageConfig() StorageConfig
	AuthConfig() AuthConfig
	UIConfig() UIConfig
	Validate() error
	GetVersionContext() string
}

type StorageConfig interface {
	GetStorageByID(storageID string) AdapterConfig
	GetStorageIDs() []string
	SigningKey() SecureString
}

type AuthConfig interface {
	GetBaseAuthConfig() *BaseAuth
	GetAuthUIConfig() *AuthUIConfig
	GetLoginURLMethodConfigParam() string
	UseUILoginPlaceholders() bool
}

type UIConfig interface {
	IsUIEnabled() bool
	GetSnippets() []apiparams.CodeSnippet
	GetCustomViewers() []apiparams.CustomViewer
}

// BaseConfig is the output struct of configuration.
type BaseConfig struct {
	ListenAddress string `mapstructure:"listen_address"`
	TLS           struct {
		Enabled  bool   `mapstructure:"enabled"`
		CertFile string `mapstructure:"cert_file"`
		KeyFile  string `mapstructure:"key_file"`
	} `mapstructure:"tls"`

	Actions struct {
		Enabled bool `mapstructure:"enabled"`
		Lua     struct {
			NetHTTPEnabled bool `mapstructure:"net_http_enabled"`
		} `mapstructure:"lua"`
		Env struct {
			Enabled bool   `mapstructure:"enabled"`
			Prefix  string `mapstructure:"prefix"`
		} `mapstructure:"env"`
	} `mapstructure:"actions"`
	Logging    Logging    `mapstructure:"logging"`
	Database   Database   `mapstructure:"database"`
	Blockstore Blockstore `mapstructure:"blockstore"`
	Committed  struct {
		LocalCache struct {
			SizeBytes             int64   `mapstructure:"size_bytes"`
			Dir                   string  `mapstructure:"dir"`
			MaxUploadersPerWriter int     `mapstructure:"max_uploaders_per_writer"`
			RangeProportion       float64 `mapstructure:"range_proportion"`
			MetaRangeProportion   float64 `mapstructure:"metarange_proportion"`
		} `mapstructure:"local_cache"`
		BlockStoragePrefix string `mapstructure:"block_storage_prefix"`
		Permanent          struct {
			MinRangeSizeBytes      uint64  `mapstructure:"min_range_size_bytes"`
			MaxRangeSizeBytes      uint64  `mapstructure:"max_range_size_bytes"`
			RangeRaggednessEntries float64 `mapstructure:"range_raggedness_entries"`
		} `mapstructure:"permanent"`
		SSTable struct {
			Memory struct {
				CacheSizeBytes int64 `mapstructure:"cache_size_bytes"`
			} `mapstructure:"memory"`
		} `mapstructure:"sstable"`
	} `mapstructure:"committed"`
	UGC struct {
		PrepareMaxFileSize int64         `mapstructure:"prepare_max_file_size"`
		PrepareInterval    time.Duration `mapstructure:"prepare_interval"`
	} `mapstructure:"ugc"`
	Graveler struct {
		EnsureReadableRootNamespace bool `mapstructure:"ensure_readable_root_namespace"`
		BatchDBIOTransactionMarkers bool `mapstructure:"batch_dbio_transaction_markers"`
		CompactionSensorThreshold   int  `mapstructure:"compaction_sensor_threshold"`
		RepositoryCache             struct {
			Size   int           `mapstructure:"size"`
			Expiry time.Duration `mapstructure:"expiry"`
			Jitter time.Duration `mapstructure:"jitter"`
		} `mapstructure:"repository_cache"`
		CommitCache struct {
			Size   int           `mapstructure:"size"`
			Expiry time.Duration `mapstructure:"expiry"`
			Jitter time.Duration `mapstructure:"jitter"`
		} `mapstructure:"commit_cache"`
		Background struct {
			RateLimit int `mapstructure:"rate_limit"`
		} `mapstructure:"background"`
		MaxBatchDelay   time.Duration                 `mapstructure:"max_batch_delay"`
		BranchOwnership ApproximatelyCorrectOwnership `mapstructure:"branch_ownership"`
	} `mapstructure:"graveler"`
	Gateways struct {
		S3 struct {
			DomainNames       Strings `mapstructure:"domain_name"`
			Region            string  `mapstructure:"region"`
			FallbackURL       string  `mapstructure:"fallback_url"`
			VerifyUnsupported bool    `mapstructure:"verify_unsupported"`
		} `mapstructure:"s3"`
	}
	Stats struct {
		Enabled       bool          `mapstructure:"enabled"`
		Address       string        `mapstructure:"address"`
		FlushInterval time.Duration `mapstructure:"flush_interval"`
		FlushSize     int           `mapstructure:"flush_size"`
		Extended      bool          `mapstructure:"extended"`
	} `mapstructure:"stats"`
	EmailSubscription struct {
		Enabled bool `mapstructure:"enabled"`
	} `mapstructure:"email_subscription"`
	Installation struct {
		FixedID                 string       `mapstructure:"fixed_id"`
		UserName                string       `mapstructure:"user_name"`
		AccessKeyID             SecureString `mapstructure:"access_key_id"`
		SecretAccessKey         SecureString `mapstructure:"secret_access_key"`
		AllowInterRegionStorage bool         `mapstructure:"allow_inter_region_storage"`
	} `mapstructure:"installation"`
	Security struct {
		CheckLatestVersion      bool          `mapstructure:"check_latest_version"`
		CheckLatestVersionCache time.Duration `mapstructure:"check_latest_version_cache"`
		AuditCheckInterval      time.Duration `mapstructure:"audit_check_interval"`
		AuditCheckURL           string        `mapstructure:"audit_check_url"`
	} `mapstructure:"security"`
	UsageReport struct {
		EnabledDeprecated bool          `mapstructure:"enabled"`
		FlushInterval     time.Duration `mapstructure:"flush_interval"`
	} `mapstructure:"usage_report"`
}

func (c *BaseConfig) GetVersionContext() string {
	return "azlake"
}

func ValidateBlockstore(c *Blockstore) error {
	if c.Signing.SecretKey == "" {
		return fmt.Errorf("'blockstore.signing.secret_key: %w", ErrMissingRequiredKeys)
	}
	if c.Type == "" {
		return fmt.Errorf("'blockstore.type: %w", ErrMissingRequiredKeys)
	}
	return nil
}

// NewConfig builds and validates a general configuration.
func NewConfig(cfgType string, c Config) (*BaseConfig, error) {
	SetDefaults(cfgType, c)
	err := Unmarshal(c)
	if err != nil {
		return nil, err
	}

	cfg := c.GetBaseConfig()
	logging.SetOutputFormat(cfg.Logging.Format)
	err = logging.SetOutputs(cfg.Logging.Output, cfg.Logging.FileMaxSizeMB, cfg.Logging.FilesKeep)
	if err != nil {
		return nil, err
	}
	logging.SetLevel(cfg.Logging.Level)
	return cfg, nil
}

func SetDefaults(cfgType string, c Config) {
	keys := GetStructKeys(reflect.TypeOf(c), "mapstructure", "squash")
	for _, key := range keys {
		viper.SetDefault(key, nil)
	}
	setBaseDefaults(cfgType)
}

func Unmarshal(c Config) error {
	return viper.UnmarshalExact(&c, DecoderConfig())
}

func DecoderConfig() viper.DecoderConfigOption {
	hook := viper.DecodeHook(
		mapstructure.ComposeDecodeHookFunc(
			DecodeStrings,
			mapstructure.StringToTimeDurationHookFunc(),
			DecodeStringToMap(),
			StringToStructHookFunc(),
			StringToSliceWithBracketHookFunc(),
		))
	return hook
}

func stringReverse(s string) string {
	chars := []rune(s)
	for i := 0; i < len(chars)/2; i++ {
		j := len(chars) - 1 - i
		chars[i], chars[j] = chars[j], chars[i]
	}
	return string(chars)
}

func (c *BaseConfig) ValidateDomainNames() error {
	domainStrings := c.Gateways.S3.DomainNames
	domainNames := make([]string, len(domainStrings))
	copy(domainNames, domainStrings)
	for i, d := range domainNames {
		domainNames[i] = stringReverse(d)
	}
	sort.Strings(domainNames)
	for i, d := range domainNames {
		domainNames[i] = stringReverse(d)
	}
	for i := 0; i < len(domainNames)-1; i++ {
		if strings.HasSuffix(domainNames[i+1], "."+domainNames[i]) {
			return fmt.Errorf("%w: %s, %s", ErrBadDomainNames, domainNames[i], domainNames[i+1])
		}
	}
	return nil
}

func (c *BaseConfig) GetBaseConfig() *BaseConfig {
	return c
}

func (c *BaseConfig) StorageConfig() StorageConfig {
	return &c.Blockstore
}

func (c *BaseConfig) Validate() error {
	missingKeys := ValidateMissingRequiredKeys(c, "mapstructure", "squash")
	if len(missingKeys) > 0 {
		return fmt.Errorf("%w: %v", ErrMissingRequiredKeys, missingKeys)
	}
	return ValidateBlockstore(&c.Blockstore)
}

type BaseAuth struct {
	Cache struct {
		Enabled bool          `mapstructure:"enabled"`
		Size    int           `mapstructure:"size"`
		TTL     time.Duration `mapstructure:"ttl"`
		Jitter  time.Duration `mapstructure:"jitter"`
	} `mapstructure:"cache"`
	Encrypt struct {
		SecretKey SecureString `mapstructure:"secret_key" validate:"required"`
	} `mapstructure:"encrypt"`
	API struct {
		Endpoint           string        `mapstructure:"endpoint"`
		Token              SecureString  `mapstructure:"token"`
		SupportsInvites    bool          `mapstructure:"supports_invites"`
		HealthCheckTimeout time.Duration `mapstructure:"health_check_timeout"`
		SkipHealthCheck    bool          `mapstructure:"skip_health_check"`
	} `mapstructure:"api"`
	AuthenticationAPI struct {
		Endpoint                  string `mapstructure:"endpoint"`
		ExternalPrincipalsEnabled bool   `mapstructure:"external_principals_enabled"`
	} `mapstructure:"authentication_api"`
	RemoteAuthenticator struct {
		Enabled          bool          `mapstructure:"enabled"`
		Endpoint         string        `mapstructure:"endpoint"`
		DefaultUserGroup string        `mapstructure:"default_user_group"`
		RequestTimeout   time.Duration `mapstructure:"request_timeout"`
	} `mapstructure:"remote_authenticator"`
	OIDC                   OIDC                   `mapstructure:"oidc"`
	CookieAuthVerification CookieAuthVerification `mapstructure:"cookie_auth_verification"`
	LogoutRedirectURL      string                 `mapstructure:"logout_redirect_url"`
	LoginDuration          time.Duration          `mapstructure:"login_duration"`
	LoginMaxDuration       time.Duration          `mapstructure:"login_max_duration"`
}

type AuthUIConfig struct {
	RBAC                 string   `mapstructure:"rbac"`
	LoginURL             string   `mapstructure:"login_url"`
	LoginFailedMessage   string   `mapstructure:"login_failed_message"`
	FallbackLoginURL     *string  `mapstructure:"fallback_login_url"`
	FallbackLoginLabel   *string  `mapstructure:"fallback_login_label"`
	LoginCookieNames     []string `mapstructure:"login_cookie_names"`
	LogoutURL            string   `mapstructure:"logout_url"`
	UseLoginPlaceholders bool     `mapstructure:"use_login_placeholders"`
}

type Auth struct {
	BaseAuth     `mapstructure:",squash"`
	AuthUIConfig `mapstructure:"ui_config"`
}

type OIDC struct {
	ValidateIDTokenClaims  map[string]string `mapstructure:"validate_id_token_claims"`
	DefaultInitialGroups   []string          `mapstructure:"default_initial_groups"`
	InitialGroupsClaimName string            `mapstructure:"initial_groups_claim_name"`
	FriendlyNameClaimName  string            `mapstructure:"friendly_name_claim_name"`
	PersistFriendlyName    bool              `mapstructure:"persist_friendly_name"`
}

type CookieAuthVerification struct {
	ValidateIDTokenClaims   map[string]string `mapstructure:"validate_id_token_claims"`
	DefaultInitialGroups    []string          `mapstructure:"default_initial_groups"`
	InitialGroupsClaimName  string            `mapstructure:"initial_groups_claim_name"`
	FriendlyNameClaimName   string            `mapstructure:"friendly_name_claim_name"`
	ExternalUserIDClaimName string            `mapstructure:"external_user_id_claim_name"`
	AuthSource              string            `mapstructure:"auth_source"`
	PersistFriendlyName     bool              `mapstructure:"persist_friendly_name"`
}

func (a *Auth) GetBaseAuthConfig() *BaseAuth {
	return &a.BaseAuth
}

func (a *Auth) GetAuthUIConfig() *AuthUIConfig {
	return &a.AuthUIConfig
}

func (a *Auth) GetLoginURLMethodConfigParam() string {
	return "none"
}

func (a *Auth) UseUILoginPlaceholders() bool {
	return a.RemoteAuthenticator.Enabled || a.AuthUIConfig.UseLoginPlaceholders
}

func (b *BaseAuth) IsAuthenticationTypeAPI() bool {
	return b.AuthenticationAPI.Endpoint != ""
}

func (b *BaseAuth) IsAuthTypeAPI() bool {
	return b.API.Endpoint != ""
}

func (b *BaseAuth) IsExternalPrincipalsEnabled() bool {
	return b.AuthenticationAPI.ExternalPrincipalsEnabled
}

func (u *AuthUIConfig) IsAuthBasic() bool {
	return u.RBAC == AuthRBACNone
}

func (u *AuthUIConfig) IsAuthUISimplified() bool {
	return u.RBAC == AuthRBACSimplified
}

func (u *AuthUIConfig) IsAdvancedAuth() bool {
	return u.RBAC == AuthRBACExternal || u.RBAC == AuthRBACInternal
}

type UI struct {
	Enabled  bool        `mapstructure:"enabled"`
	Snippets []UISnippet `mapstructure:"snippets"`
}

type UISnippet struct {
	ID   string `mapstructure:"id"`
	Code string `mapstructure:"code"`
}

func (u *UI) IsUIEnabled() bool {
	return u.Enabled
}

func (u *UI) GetSnippets() []apiparams.CodeSnippet {
	return BuildCodeSnippets(u.Snippets)
}

func BuildCodeSnippets(s []UISnippet) []apiparams.CodeSnippet {
	snippets := make([]apiparams.CodeSnippet, 0, len(s))
	for _, item := range s {
		snippets = append(snippets, apiparams.CodeSnippet{
			ID:   item.ID,
			Code: item.Code,
		})
	}
	return snippets
}

func (u *UI) GetCustomViewers() []apiparams.CustomViewer {
	return nil
}
