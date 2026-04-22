package config

import (
	"time"

	"github.com/spf13/viper"
)

const (
	DefaultListenAddress        = "0.0.0.0:8000"
	DefaultLoggingLevel         = "INFO"
	DefaultLoggingAuditLogLevel = "DEBUG"
	DefaultLoggingFilesKeep     = 100
	DefaultLoggingFileMaxSizeMB = 1024 * 100

	BlockstoreTypeKey                        = "blockstore.type"
	DefaultQuickstartUsername                 = "quickstart"
	DefaultQuickstartKeyID                   = "AZLAKE_QUICKSTART_KEY_ID"     //nolint:gosec
	DefaultQuickstartSecretKey               = "AZLAKE_QUICKSTART_SECRET_KEY" //nolint:gosec
	DefaultAuthAPIHealthCheckTimeout         = 20 * time.Second
	DefaultAuthSecret                        = "THIS_MUST_BE_CHANGED_IN_PRODUCTION"   // #nosec
	DefaultSigningSecretKey                  = "OVERRIDE_THIS_SIGNING_SECRET_DEFAULT" // #nosec
	DefaultBlockstoreLocalPath               = "~/lakefs/data/block"
	DefaultBlockstoreAzureTryTimeout         = 10 * time.Minute
	DefaultBlockstoreAzurePreSignedExpiry    = 15 * time.Minute
	DefaultBlockstoreAzureDisablePreSignedUI = true
)

//nolint:mnd
func setBaseDefaults(cfgType string) {
	switch cfgType {
	case QuickstartConfiguration:
		viper.SetDefault("installation.user_name", DefaultQuickstartUsername)
		viper.SetDefault("installation.access_key_id", DefaultQuickstartKeyID)
		viper.SetDefault("installation.secret_access_key", DefaultQuickstartSecretKey)
		viper.SetDefault("database.type", "local")
		viper.SetDefault("auth.encrypt.secret_key", DefaultAuthSecret)
		viper.SetDefault(BlockstoreTypeKey, "local")
	case UseLocalConfiguration:
		viper.SetDefault("database.type", "local")
		viper.SetDefault("auth.encrypt.secret_key", DefaultAuthSecret)
		viper.SetDefault(BlockstoreTypeKey, "local")
	}

	viper.SetDefault("installation.allow_inter_region_storage", true)
	viper.SetDefault("listen_address", DefaultListenAddress)

	SetLoggingDefaults()

	viper.SetDefault("actions.enabled", true)
	viper.SetDefault("actions.env.enabled", true)
	viper.SetDefault("actions.env.prefix", "LAKEFSACTION_")

	viper.SetDefault("auth.cache.enabled", true)
	viper.SetDefault("auth.cache.size", 1024)
	viper.SetDefault("auth.cache.ttl", 20*time.Second)
	viper.SetDefault("auth.cache.jitter", 3*time.Second)

	viper.SetDefault("auth.logout_redirect_url", "/auth/login")
	viper.SetDefault("auth.login_duration", 7*24*time.Hour)
	viper.SetDefault("auth.login_max_duration", 14*24*time.Hour)

	viper.SetDefault("auth.ui_config.rbac", "none")
	viper.SetDefault("auth.ui_config.login_failed_message", "The credentials don't match.")
	viper.SetDefault("auth.ui_config.login_cookie_names", "internal_auth_session")
	viper.SetDefault("auth.ui_config.use_login_placeholders", false)

	viper.SetDefault("auth.remote_authenticator.default_user_group", "Viewers")
	viper.SetDefault("auth.remote_authenticator.request_timeout", 10*time.Second)

	viper.SetDefault("auth.api.health_check_timeout", DefaultAuthAPIHealthCheckTimeout)
	viper.SetDefault("auth.oidc.persist_friendly_name", false)
	viper.SetDefault("auth.cookie_auth_verification.persist_friendly_name", false)

	viper.SetDefault("committed.local_cache.size_bytes", 1*1024*1024*1024)
	viper.SetDefault("committed.local_cache.dir", "~/lakefs/data/cache")
	viper.SetDefault("committed.local_cache.max_uploaders_per_writer", 10)
	viper.SetDefault("committed.local_cache.range_proportion", 0.9)
	viper.SetDefault("committed.local_cache.metarange_proportion", 0.1)

	viper.SetDefault("committed.block_storage_prefix", "_lakefs")

	viper.SetDefault("committed.permanent.min_range_size_bytes", 0)
	viper.SetDefault("committed.permanent.max_range_size_bytes", 20*1024*1024)
	viper.SetDefault("committed.permanent.range_raggedness_entries", 50_000)

	viper.SetDefault("committed.sstable.memory.cache_size_bytes", 400_000_000)

	viper.SetDefault("gateways.s3.domain_name", "s3.local.lakefs.io")
	viper.SetDefault("gateways.s3.region", "us-east-1")
	viper.SetDefault("gateways.s3.verify_unsupported", true)

	// blockstore defaults
	viper.SetDefault("blockstore.signing.secret_key", DefaultSigningSecretKey)
	viper.SetDefault("blockstore.local.path", DefaultBlockstoreLocalPath)

	viper.SetDefault("blockstore.azure.try_timeout", DefaultBlockstoreAzureTryTimeout)
	viper.SetDefault("blockstore.azure.pre_signed_expiry", DefaultBlockstoreAzurePreSignedExpiry)
	viper.SetDefault("blockstore.azure.disable_pre_signed_ui", DefaultBlockstoreAzureDisablePreSignedUI)

	viper.SetDefault("stats.enabled", true)
	viper.SetDefault("stats.address", "https://stats.lakefs.io")
	viper.SetDefault("stats.flush_interval", 30*time.Second)
	viper.SetDefault("stats.flush_size", 100)

	viper.SetDefault("email_subscription.enabled", true)

	viper.SetDefault("security.audit_check_interval", 24*time.Hour)
	viper.SetDefault("security.audit_check_url", "https://audit.lakefs.io/audit")
	viper.SetDefault("security.check_latest_version", true)
	viper.SetDefault("security.check_latest_version_cache", time.Hour)

	viper.SetDefault("ui.enabled", true)

	viper.SetDefault("database.local.path", "~/lakefs/metadata")
	viper.SetDefault("database.local.prefetch_size", 256)
	viper.SetDefault("database.local.sync_writes", true)

	viper.SetDefault("graveler.ensure_readable_root_namespace", true)
	viper.SetDefault("graveler.repository_cache.size", 1000)
	viper.SetDefault("graveler.repository_cache.expiry", 5*time.Second)
	viper.SetDefault("graveler.repository_cache.jitter", 2*time.Second)
	viper.SetDefault("graveler.commit_cache.size", 50_000)
	viper.SetDefault("graveler.commit_cache.expiry", 10*time.Minute)
	viper.SetDefault("graveler.commit_cache.jitter", 2*time.Second)

	viper.SetDefault("graveler.max_batch_delay", 3*time.Millisecond)

	viper.SetDefault("graveler.branch_ownership.enabled", false)
	viper.SetDefault("graveler.branch_ownership.refresh", 400*time.Millisecond)
	viper.SetDefault("graveler.branch_ownership.acquire", 150*time.Millisecond)

	viper.SetDefault("ugc.prepare_interval", time.Minute)
	viper.SetDefault("ugc.prepare_max_file_size", 20*1024*1024)

	viper.SetDefault("usage_report.flush_interval", 5*time.Minute)
}

func SetLoggingDefaults() {
	viper.SetDefault("logging.format", "text")
	viper.SetDefault("logging.level", DefaultLoggingLevel)
	viper.SetDefault("logging.output", "-")
	viper.SetDefault("logging.files_keep", DefaultLoggingFilesKeep)
	viper.SetDefault("logging.audit_log_level", DefaultLoggingAuditLogLevel)
	viper.SetDefault("logging.file_max_size_mb", DefaultLoggingFileMaxSizeMB)
}
