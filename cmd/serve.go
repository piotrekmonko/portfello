/*
Copyright © 2023 Piotr Mońko <piotrek.monko@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"context"
	"errors"
	"github.com/piotrekmonko/portfello/pkg/auth"
	"github.com/piotrekmonko/portfello/pkg/config"
	"github.com/piotrekmonko/portfello/pkg/dao"
	"github.com/piotrekmonko/portfello/pkg/server"
	"github.com/spf13/cobra"
	"log"
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
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		conf := config.New()

		db, dbQuerier, err := dao.NewDAO(ctx, conf.DatabaseDSN)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		authProvider, err := auth.NewAuth0Provider(cmd.Context(), conf)
		if err != nil {
			log.Fatal(err)
		}
		authService := auth.New(authProvider)

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

		httpSrv := server.NewServer(ctx, conf, dbQuerier, authService)
		go func() {
			if err := httpSrv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
				log.Fatal(err)
			}
		}()
		<-sigs

		closeCtx, closeCanc := context.WithTimeout(ctx, time.Second)
		defer closeCanc()
		cobra.CheckErr(httpSrv.Shutdown(closeCtx))
		log.Println("Server stopped")
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
