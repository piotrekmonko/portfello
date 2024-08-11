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

func NewServer(ctx context.Context, log *logz.Log, conf *conf.Config, mux *http.ServeMux) (*http.Server, func(), error) {

	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:" + conf.Graph.Port, "http://" + conf.HostName},
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
	return httpSrv, Shutdown(log, httpSrv), nil
}

func Shutdown(log logz.Logger, s *http.Server) func() {
	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		log.Infow(ctx, "shutting down server")
		_ = log.Errorw(ctx, s.Shutdown(ctx), "error shutting down server")
	}
}
