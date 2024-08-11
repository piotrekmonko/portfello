package cmd

import (
	"errors"
	"github.com/fsnotify/fsnotify"
	"github.com/piotrekmonko/portfello/pkg/conf"
	"github.com/piotrekmonko/portfello/pkg/logz"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start GraphQL server",
	RunE: func(cmd *cobra.Command, _ []string) error {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGUSR1)

		viper.OnConfigChange(func(_ fsnotify.Event) {
			log.Println("reloading server due to config change...")
			sigs <- syscall.SIGUSR1
		})

		return serveWithRestartOnConfigChange(cmd, sigs)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func serveWithRestartOnConfigChange(cmd *cobra.Command, sigs chan os.Signal) error {
	ctx := cmd.Context()
	c := conf.New()

	log, syncer, err := logz.NewLogger(c)
	if err != nil {
		return err
	}

	httpSrv, httpCloser, err := initializeServer(cmd.Context(), c)
	if err != nil {
		return err
	}

	go func() {
		if errz := httpSrv.ListenAndServe(); !errors.Is(errz, http.ErrServerClosed) {
			_ = log.Errorw(ctx, errz, "error while running http server")
			sigs <- syscall.SIGQUIT
		}
	}()

	reason := <-sigs
	httpCloser()
	log.Infow(ctx, "Stopped")
	syncer()

	if reason == syscall.SIGUSR1 {
		return serveWithRestartOnConfigChange(cmd, sigs)
	}

	return nil
}
