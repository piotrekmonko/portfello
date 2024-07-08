package cmd

import (
	"context"
	"errors"
	"github.com/piotrekmonko/portfello/pkg/conf"
	"github.com/piotrekmonko/portfello/pkg/logz"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start GraphQL server",
	RunE: func(cmd *cobra.Command, _ []string) error {
		ctx := cmd.Context()
		c := conf.New()

		log, syncer, err := logz.NewLogger(c)
		if err != nil {
			return err
		}
		defer syncer()

		httpSrv, closer, err := initializeServer(cmd.Context(), c)
		if err != nil {
			return err
		}
		defer closer()

		go func() {
			if err := httpSrv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
				_ = log.Errorw(ctx, err, "error while running http server")
			}
		}()

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		<-sigs

		log.Infow(ctx, "Stopping server...")
		closeCtx, closeCanc := context.WithTimeout(ctx, time.Second)
		defer closeCanc()
		cobra.CheckErr(httpSrv.Shutdown(closeCtx))
		log.Infow(ctx, "Stopped")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
