package server

import (
	"context"
	"github.com/piotrekmonko/portfello/pkg/conf"
	"github.com/piotrekmonko/portfello/pkg/logz"
	"github.com/rs/cors"
	"go.uber.org/zap/zapcore"
	"net/http"
	"time"
)

func NewServer(ctx context.Context, log *logz.Log, conf *conf.Config, mux *http.ServeMux) (*http.Server, error) {

	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"localhost:" + conf.Graph.Port, conf.HostName},
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		MaxAge: 30,
		Debug:  log.Level() == zapcore.DebugLevel,
		Logger: log,
	})

	httpSrv := &http.Server{
		Addr:              ":" + conf.Graph.Port,
		Handler:           corsHandler.Handler(mux),
		ReadHeaderTimeout: time.Second,
	}

	log.Infof(ctx, "serving on http://localhost:%s/", conf.Graph.Port)
	return httpSrv, nil
}
