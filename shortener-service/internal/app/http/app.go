package httpapp

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/nikitamavrenko/shortener-service/internal/config"
	"github.com/nikitamavrenko/shortener-service/internal/http-routes/redirect"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"net/http"
)

type App struct {
	server *http.Server
	log    *zerolog.Logger
}

func New(cfg *config.Config, log *zerolog.Logger, redirector redirect.Redirector) *App {
	router := mux.NewRouter()

	router.HandleFunc("/{id}", redirect.New(log, redirector)).Methods(http.MethodGet)

	router.Handle("/metrics", promhttp.Handler()).Methods(http.MethodGet)

	server := &http.Server{
		Handler: router,
		Addr:    ":8080",
	}

	return &App{server: server, log: log}
}

func (a *App) Run() error {
	a.log.Info().Msg("running web server")
	err := a.server.ListenAndServe()
	if err != nil {
		a.log.Warn().Err(err)
		return err
	}
	a.log.Info().Msg("web server is started on address " + a.server.Addr)
	return nil
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		a.log.Fatal().Err(err).Msg("shortener server failed")
		panic(err)
	}
}

func (a *App) Shutdown(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}
