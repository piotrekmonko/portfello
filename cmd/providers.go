//go:build wireinject
// +build wireinject

package cmd

import (
	"context"
	"github.com/google/wire"
	"github.com/piotrekmonko/portfello/pkg/auth"
	"github.com/piotrekmonko/portfello/pkg/conf"
	"github.com/piotrekmonko/portfello/pkg/dao"
	"github.com/piotrekmonko/portfello/pkg/logz"
	"github.com/piotrekmonko/portfello/pkg/provision"
	"github.com/piotrekmonko/portfello/pkg/server"
	"net/http"
)

func initializeProvisioner(ctx context.Context, c *conf.Config) (*provision.Provisioner, func(), error) {
	wire.Build(provision.NewProvisioner, dao.NewDAO, auth.NewFromConfig, logz.NewLogger)
	return &provision.Provisioner{}, func() {}, nil
}

func initializeServer(ctx context.Context, c *conf.Config) (*http.Server, func(), error) {
	wire.Build(server.NewServer, server.NewRouter, dao.NewDAO, auth.NewFromConfig, logz.NewLogger)
	return &http.Server{}, nil, nil
}

func initializeRouter(ctx context.Context, c *conf.Config) (*http.ServeMux, func(), error) {
	wire.Build(server.NewRouter, dao.NewDAO, auth.NewFromConfig, logz.NewLogger)
	return &http.ServeMux{}, nil, nil
}
