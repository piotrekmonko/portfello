package server

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/hellofresh/health-go/v5"
	"github.com/piotrekmonko/portfello/pkg/auth"
	"github.com/piotrekmonko/portfello/pkg/conf"
	"github.com/piotrekmonko/portfello/pkg/dao"
	"github.com/piotrekmonko/portfello/pkg/graph"
	"github.com/piotrekmonko/portfello/pkg/logz"
	"net/http"
	"time"
)

func NewServer(ctx context.Context, log logz.Logger, conf *conf.Config, dbQuerier *dao.DAO, authService *auth.Service) *http.Server {
	graphResolver := &graph.Resolver{
		Conf:        conf,
		DbDAO:       dbQuerier,
		AuthService: authService,
	}

	graphConfig := graph.Config{
		Resolvers: graphResolver,
	}
	graphConfig.Directives.HasRole = authService.HasRole

	mux := http.NewServeMux()
	httpSrv := &http.Server{
		Addr:              ":" + conf.Graph.Port,
		Handler:           mux,
		ReadHeaderTimeout: time.Second,
	}

	healthChecks, err := health.New(
		health.WithComponent(health.Component{
			Name:    "portfello",
			Version: logz.GetVer(),
		}),
		health.WithSystemInfo(),
		health.WithChecks(
			health.Config{
				Name: "postgresql",
				Check: func(ctx context.Context) error {
					return dbQuerier.Ping(ctx)
				},
			},
		),
	)
	if err != nil {
		panic(fmt.Errorf("cannot start healthcheck service: %w", err))
	}

	mux.Handle("/healthcheck", healthChecks.Handler())

	srv := authService.Middleware(handler.NewDefaultServer(graph.NewExecutableSchema(graphConfig)))
	mux.Handle("/query", srv)

	if conf.Graph.EnablePlayground {
		log.Infof(ctx, "connect to http://localhost:%s/ for GraphQL playground", conf.Graph.Port)
		mux.Handle("/", playground.Handler("GraphQL playground", "/query"))
	}

	log.Infof(ctx, "serving on http://localhost:%s/", conf.Graph.Port)
	return httpSrv
}
