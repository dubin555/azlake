package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/dubin555/azlake/pkg/api"
	"github.com/dubin555/azlake/pkg/azcat"
	"github.com/dubin555/azlake/pkg/logging"
	"github.com/dubin555/azlake/pkg/version"
)

const gracefulShutdownTimeout = 30 * time.Second

func getEnvOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the azlake server",
	Run: func(cmd *cobra.Command, args []string) {
		logger := logging.ContextUnavailable()
		logger.WithField("version", version.Version).Info("Starting azlake")

		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		// ── KV backend ──
		kvBackend := getEnvOrDefault("AZLAKE_KV_BACKEND", "badger")
		var kv azcat.KV

		switch kvBackend {
		case "cosmosdb":
			endpoint := os.Getenv("COSMOS_ENDPOINT")
			key := os.Getenv("COSMOS_KEY")
			database := getEnvOrDefault("COSMOS_DATABASE", "azlake")
			container := getEnvOrDefault("COSMOS_CONTAINER", "metadata")
			logger.WithField("endpoint", endpoint).WithField("database", database).Info("Connecting to CosmosDB")
			store, err := azcat.OpenCosmosKV(endpoint, key, database, container)
			if err != nil {
				logger.WithError(err).Fatal("Failed to open CosmosDB KV store")
			}
			kv = store
		default: // "badger"
			dataDir := azcat.DefaultDataDir()
			logger.WithField("data_dir", dataDir).Info("Opening BadgerDB")
			store, err := azcat.OpenKV(dataDir)
			if err != nil {
				logger.WithError(err).Fatal("Failed to open KV store")
			}
			kv = store
		}
		defer kv.Close()

		// ── Object storage backend ──
		storageBackend := getEnvOrDefault("AZLAKE_STORAGE_BACKEND", "local")
		var storage azcat.ObjectStorage

		switch storageBackend {
		case "azure":
			accountName := os.Getenv("AZURE_STORAGE_ACCOUNT")
			accountKey := os.Getenv("AZURE_STORAGE_KEY") // optional, falls back to DefaultAzureCredential
			containerName := getEnvOrDefault("AZURE_STORAGE_CONTAINER", "azlake")
			logger.WithField("account", accountName).WithField("container", containerName).Info("Using Azure Blob Storage")
			azStorage, err := azcat.NewAzureBlobStorage(accountName, accountKey, containerName)
			if err != nil {
				logger.WithError(err).Fatal("Failed to initialize Azure Blob Storage")
			}
			storage = azStorage
		default: // "local"
			objectsDir := azcat.DefaultObjectsDir()
			logger.WithField("objects_dir", objectsDir).Info("Using local object storage")
			storage = azcat.NewLocalStorage(objectsDir)
		}

		catalog := azcat.NewCatalog(kv, storage)

		listenAddr := ":8000"
		logger.WithField("listen_address", listenAddr).Info("Starting HTTP server")

		handler := api.Serve(logger, catalog)

		server := &http.Server{
			Addr:    listenAddr,
			Handler: handler,
		}

		go func() {
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logger.WithError(err).Fatal("HTTP server failed")
			}
		}()

		<-ctx.Done()
		logger.Info("Shutting down...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimeout)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.WithError(err).Error("server shutdown failed")
		}
		fmt.Println("azlake stopped")
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
