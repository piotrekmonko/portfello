package server

import (
	"context"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/piotrekmonko/portfello/pkg/auth"
	"github.com/piotrekmonko/portfello/pkg/config"
	"github.com/piotrekmonko/portfello/pkg/dao"
	"github.com/piotrekmonko/portfello/pkg/graph"
	"log"
	"net/http"
	"time"
)

func NewServer(ctx context.Context, conf *config.Config, dbQuerier *dao.DAO, authService *auth.Service) *http.Server {
	graphResolver := &graph.Resolver{
		Conf:        conf,
		DbDAO:       dbQuerier,
		AuthService: authService,
	}

	graphConfig := graph.Config{
		Resolvers: graphResolver,
	}
	graphConfig.Directives.HasRole = authService.HasRole

	srv := authService.Middleware(handler.NewDefaultServer(graph.NewExecutableSchema(graphConfig)))
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
	return httpSrv
}