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
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/piotrekmonko/portfello/pkg/auth"
	"github.com/piotrekmonko/portfello/pkg/config"
	"github.com/piotrekmonko/portfello/pkg/dao"
	"github.com/piotrekmonko/portfello/pkg/graph"
	"log"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start GraphQL server",
	Run: func(cmd *cobra.Command, args []string) {
		conf := config.New()

		db, dbQuerier, err := dao.NewDAO(conf.DatabaseDSN)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		authProvider, err := auth.NewAuth0Provider(cmd.Context(), conf)
		if err != nil {
			log.Fatal(err)
		}
		authService := auth.New(authProvider)

		graphResolver := &graph.Resolver{
			Conf:        conf,
			DbQueries:   dbQuerier,
			AuthService: authService,
		}
		srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: graphResolver}))
		mux := http.NewServeMux()
		httpSrv := &http.Server{
			Addr:              ":" + conf.Graph.Port,
			Handler:           mux,
			ReadHeaderTimeout: time.Second,
		}

		mux.Handle("/query", srv)
		if conf.Graph.EnablePlayground {
			log.Printf("connect to http://localhost:%s/ for GraphQL playground", conf.Graph.Port)
			mux.Handle("/", playground.Handler("GraphQL playground", "/query"))
		}

		log.Printf("serving on http://localhost:%s/", conf.Graph.Port)
		log.Fatal(httpSrv.ListenAndServe())
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
