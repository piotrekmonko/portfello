package server

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/hellofresh/health-go/v5"
	"github.com/piotrekmonko/portfello/pkg/auth"
	"github.com/piotrekmonko/portfello/pkg/conf"
	"github.com/piotrekmonko/portfello/pkg/dao"
	"github.com/piotrekmonko/portfello/pkg/graph"
	"github.com/piotrekmonko/portfello/pkg/logz"
	"net/http"
)

// NewRouter builds routing mux, registers handlers and health checks.
func NewRouter(conf *conf.Config, dbQuerier *dao.DAO, authService *auth.Service) (*http.ServeMux, error) {
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
		return nil, fmt.Errorf("cannot start healthcheck service: %w", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/log/level", logz.AtomicLevel)
	mux.Handle("/healthcheck", healthChecks.Handler())
	mux.Handle("/query", graph.NewGraphHandler(conf, dbQuerier, authService))

	if conf.Graph.EnablePlayground {
		mux.Handle("/", playground.Handler("GraphQL playground", "/query"))
	}

	return mux, nil
}
