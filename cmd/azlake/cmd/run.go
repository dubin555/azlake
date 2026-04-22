package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/dubin555/azlake/pkg/api"
	"github.com/dubin555/azlake/pkg/logging"
	"github.com/dubin555/azlake/pkg/version"
)

const gracefulShutdownTimeout = 30 * time.Second

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the azlake server",
	Run: func(cmd *cobra.Command, args []string) {
		logger := logging.ContextUnavailable()
		logger.WithField("version", version.Version).Info("Starting azlake")

		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		listenAddr := ":8000"
		logger.WithField("listen_address", listenAddr).Info("Starting HTTP server")

		handler := api.Serve(logger)

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
