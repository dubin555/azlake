package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/deepmap/oapi-codegen/pkg/securityprovider"
	"github.com/go-openapi/swag"
	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/dubin555/azlake/pkg/api/apigen"
	"github.com/dubin555/azlake/pkg/api/apiutil"
	lakefsconfig "github.com/dubin555/azlake/pkg/config"
	"github.com/dubin555/azlake/pkg/logging"
	"github.com/dubin555/azlake/pkg/version"
)

const (
	DefaultMaxIdleConnsPerHost = 100

	versionTemplate = `azlakectl version: {{.Version }}
{{- if .ServerVersion }}
azlake server version: {{.ServerVersion}}
{{- end }}
`
)

// Configuration is the CLI configuration structure.
type Configuration struct {
	Credentials struct {
		AccessKeyID     lakefsconfig.OnlyString `mapstructure:"access_key_id"`
		SecretAccessKey  lakefsconfig.OnlyString `mapstructure:"secret_access_key"`
	} `mapstructure:"credentials"`
	Server struct {
		EndpointURL lakefsconfig.OnlyString `mapstructure:"endpoint_url"`
	} `mapstructure:"server"`
}

type versionInfo struct {
	Version       string
	ServerVersion string
}

var (
	cfgFile string
	cfgErr  error
	cfg     *Configuration

	baseURI  string
	logLevel string
	logFormat string
	logOutputs []string

	noColorRequested = false
	verboseMode      = false
)

const (
	recursiveFlagName     = "recursive"
	recursiveFlagShort    = "r"
	presignFlagName       = "pre-sign"
	parallelismFlagName   = "parallelism"

	defaultParallelism = 25
	defaultSyncPresign = true

	paginationPrefixFlagName = "prefix"
	paginationAfterFlagName  = "after"
	paginationAmountFlagName = "amount"

	myRepoExample   = "azlake://my-repo"
	myBranchExample = "my-branch"
	myBucketExample = "az://my-container"
	myDigestExample = "600dc0ffee"
	myRunIDExample  = "20230719152411arS0z6I"
	storageIDFlagName = "storage-id"
	noProgressBarFlagName = "no-progress"
	defaultNoProgress  = false

	commitMsgFlagName     = "message"
	allowEmptyMsgFlagName = "allow-empty-message"
	fmtErrEmptyMsg        = `commit with no message without specifying the "--allow-empty-message" flag`
	metaFlagName          = "meta"
)

func withRecursiveFlag(cmd *cobra.Command, usage string) {
	cmd.Flags().BoolP(recursiveFlagName, recursiveFlagShort, false, usage)
}

func withParallelismFlag(cmd *cobra.Command) {
	cmd.Flags().IntP(parallelismFlagName, "p", defaultParallelism,
		"Max concurrent operations to perform")
}

func withPresignFlag(cmd *cobra.Command) {
	cmd.Flags().Bool(presignFlagName, defaultSyncPresign,
		"Use pre-signed URLs when downloading/uploading data (recommended)")
}

func withStorageID(cmd *cobra.Command) {
	cmd.Flags().String(storageIDFlagName, "", "")
	if err := cmd.Flags().MarkHidden(storageIDFlagName); err != nil {
		DieErr(err)
	}
}

func withNoProgress(cmd *cobra.Command) {
	cmd.Flags().Bool(noProgressBarFlagName, defaultNoProgress,
		"Disable progress bar animation for IO operations")
}

func withSyncFlags(cmd *cobra.Command) {
	withParallelismFlag(cmd)
	withPresignFlag(cmd)
	withNoProgress(cmd)
}

func getPresignMode(cmd *cobra.Command, client *apigen.ClientWithResponses, repositoryID string) PresignMode {
	presignFlag := cmd.Flags().Lookup(presignFlagName)
	if presignFlag != nil && presignFlag.Changed {
		return PresignMode{Enabled: Must(cmd.Flags().GetBool(presignFlagName))}
	}
	return getServerPreSignMode(cmd.Context(), client, repositoryID)
}

func getPaginationFlags(cmd *cobra.Command) (prefix string, after string, amount int) {
	prefix = Must(cmd.Flags().GetString(paginationPrefixFlagName))
	after = Must(cmd.Flags().GetString(paginationAfterFlagName))
	amount = Must(cmd.Flags().GetInt(paginationAmountFlagName))
	return
}

type PaginationOptions func(*cobra.Command)

func withoutPrefix(cmd *cobra.Command) {
	if err := cmd.Flags().MarkHidden(paginationPrefixFlagName); err != nil {
		DieErr(err)
	}
}

func withPaginationFlags(cmd *cobra.Command, options ...PaginationOptions) {
	cmd.Flags().SortFlags = false
	cmd.Flags().Int(paginationAmountFlagName, defaultAmountArgumentValue, "how many results to return")
	cmd.Flags().String(paginationAfterFlagName, "", "show results after this value (used for pagination)")
	cmd.Flags().String(paginationPrefixFlagName, "", "filter results by prefix (used for pagination)")
	for _, option := range options {
		option(cmd)
	}
}

func withMessageFlags(cmd *cobra.Command, allowEmpty bool) {
	cmd.Flags().StringP(commitMsgFlagName, "m", "", "commit message")
	cmd.Flags().Bool(allowEmptyMsgFlagName, allowEmpty, "allow an empty commit message")
}

func withMetadataFlag(cmd *cobra.Command) {
	cmd.Flags().StringSlice(metaFlagName, []string{}, "key value pair in the form of key=value")
}

func withCommitFlags(cmd *cobra.Command, allowEmptyMessage bool) {
	withMessageFlags(cmd, allowEmptyMessage)
	withMetadataFlag(cmd)
}

func getCommitFlags(cmd *cobra.Command) (string, map[string]string) {
	message := Must(cmd.Flags().GetString(commitMsgFlagName))
	emptyMessageBool := Must(cmd.Flags().GetBool(allowEmptyMsgFlagName))
	if strings.TrimSpace(message) == "" && !emptyMessageBool {
		DieFmt(fmtErrEmptyMsg)
	}
	kvPairs, err := getKV(cmd, metaFlagName)
	if err != nil {
		DieErr(err)
	}
	return message, kvPairs
}

func getKV(cmd *cobra.Command, name string) (map[string]string, error) {
	kvList, err := cmd.Flags().GetStringSlice(name)
	if err != nil {
		return nil, err
	}
	kv := make(map[string]string)
	for _, pair := range kvList {
		key, value, found := strings.Cut(pair, "=")
		if !found {
			return nil, errInvalidKeyValueFormat
		}
		kv[key] = value
	}
	return kv, nil
}

// rootCmd represents the base command when called without any sub-commands
var rootCmd = &cobra.Command{
	Use:   "azlakectl",
	Short: "Azure-native data version control CLI",
	Long:  `azlakectl is a CLI tool for managing azlake repositories, branches, commits, and data.`,
	SilenceErrors: true,
	SilenceUsage:  true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logging.SetLevel(logLevel)
		logging.SetOutputFormat(logFormat)
		if err := logging.SetOutputs(logOutputs, 0, 0); err != nil {
			DieFmt("Failed to setup logging: %s", err)
		}
		if noColorRequested {
			DisableColors()
		}
		if cmd == configCmd {
			return
		}
		if cfgFile != "" && cfgErr != nil {
			DieFmt("error reading configuration file: %v", cfgErr)
		}
		if err := viper.UnmarshalExact(&cfg, viper.DecodeHook(
			mapstructure.ComposeDecodeHookFunc(
				lakefsconfig.DecodeOnlyString,
				mapstructure.StringToTimeDurationHookFunc(),
				lakefsconfig.DecodeStringToMap(),
			))); err != nil {
			DieFmt("error unmarshal configuration: %v", err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if !Must(cmd.Flags().GetBool("version")) {
			if err := cmd.Help(); err != nil {
				WriteIfVerbose("failed showing help {{ . }}", err)
			}
			return
		}
		info := versionInfo{Version: version.Version}
		client := getClient()
		resp, err := client.GetConfigWithResponse(cmd.Context())
		if err == nil && resp.JSON200 != nil {
			info.ServerVersion = swag.StringValue(resp.JSON200.VersionConfig.Version)
		}
		Write(versionTemplate, info)
	},
}

func getHTTPClient() *http.Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConnsPerHost = DefaultMaxIdleConnsPerHost
	return &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}
}

func getClient(opts ...apigen.ClientOption) *apigen.ClientWithResponses {
	httpClient := getHTTPClient()
	accessKeyID := cfg.Credentials.AccessKeyID
	secretAccessKey := cfg.Credentials.SecretAccessKey
	basicAuthProvider, err := securityprovider.NewSecurityProviderBasicAuth(string(accessKeyID), string(secretAccessKey))
	if err != nil {
		DieErr(err)
	}
	serverEndpoint, err := apiutil.NormalizeLakeFSEndpoint(cfg.Server.EndpointURL.String())
	if err != nil {
		DieErr(err)
	}
	client, err := apigen.NewClientWithResponses(
		serverEndpoint,
		append([]apigen.ClientOption{
			apigen.WithHTTPClient(httpClient),
			apigen.WithRequestEditorFn(basicAuthProvider.Intercept),
			apigen.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
				req.Header.Set("User-Agent", fmt.Sprintf("azlakectl/%s", version.Version))
				return nil
			}),
		}, opts...)...,
	)
	if err != nil {
		Die(fmt.Sprintf("could not initialize API client: %s", err), 1)
	}
	return client
}

func getStorageConfigOrDie(ctx context.Context, client *apigen.ClientWithResponses, repositoryID string) *apigen.StorageConfig {
	confResp, err := client.GetConfigWithResponse(ctx)
	DieOnErrorOrUnexpectedStatusCode(confResp, err, http.StatusOK)
	if confResp.JSON200 == nil {
		Die("Bad response from server for GetConfig", 1)
	}
	storageConfig := confResp.JSON200.StorageConfig
	if storageConfig == nil {
		Die("Bad response from server for GetConfig", 1)
	}
	return storageConfig
}

type PresignMode struct {
	Enabled   bool
	Multipart bool
}

func getServerPreSignMode(ctx context.Context, client *apigen.ClientWithResponses, repositoryID string) PresignMode {
	storageConfig := getStorageConfigOrDie(ctx, client, repositoryID)
	return PresignMode{
		Enabled:   storageConfig.PreSignSupport,
		Multipart: swag.BoolValue(storageConfig.PreSignMultipartUpload),
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

//nolint:gochecknoinits
func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.azlakectl.yaml)")
	rootCmd.PersistentFlags().BoolVar(&noColorRequested, "no-color", getEnvNoColor(), "don't use fancy output colors")
	rootCmd.PersistentFlags().StringVarP(&baseURI, "base-uri", "", os.Getenv("AZLAKECTL_BASE_URI"), "base URI used for azlake address parse")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "", "none", "set logging level")
	rootCmd.PersistentFlags().StringVarP(&logFormat, "log-format", "", "", "set logging output format")
	rootCmd.PersistentFlags().StringSliceVarP(&logOutputs, "log-output", "", []string{}, "set logging output(s)")
	rootCmd.PersistentFlags().BoolVar(&verboseMode, "verbose", false, "run in verbose mode")
	rootCmd.Flags().BoolP("version", "v", false, "version for azlakectl")
}

func getEnvNoColor() bool {
	v := os.Getenv("NO_COLOR")
	return v != "" && v != "0"
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			DieErr(err)
		}
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".azlakectl")
	}
	viper.SetEnvPrefix("AZLAKECTL")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	keys := lakefsconfig.GetStructKeys(reflect.TypeFor[Configuration](), "mapstructure", "squash")
	for _, key := range keys {
		viper.SetDefault(key, nil)
	}
	cfgErr = viper.ReadInConfig()
}
