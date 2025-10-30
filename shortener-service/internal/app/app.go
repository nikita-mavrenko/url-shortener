package app

import (
	"context"
	"fmt"
	grpcapp "github.com/nikitamavrenko/shortener-service/internal/app/grpc"
	httpapp "github.com/nikitamavrenko/shortener-service/internal/app/http"
	"github.com/nikitamavrenko/shortener-service/internal/config"
	"github.com/nikitamavrenko/shortener-service/internal/services/shortener"
	"github.com/nikitamavrenko/shortener-service/internal/storage/postgres"
	"github.com/rs/zerolog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	log     *zerolog.Logger
	gRPCApp *grpcapp.App
	httpApp *httpapp.App
}

func New(log *zerolog.Logger, cfg *config.Config) *App {
	db, err := postgres.New(context.Background(), cfg.Db.Url, log)
	if err != nil {
		panic(err)
	}
	err = db.Migrate(context.Background())
	if err != nil {
		panic(err)
	}

	shortenerService := shortener.New(log, db, cfg.Http.Address)

	gRPCApp := grpcapp.New(log, shortenerService, cfg.GRPC.Port)

	httpApp := httpapp.New(cfg, log, shortenerService)

	return &App{
		log:     log,
		gRPCApp: gRPCApp,
		httpApp: httpApp,
	}
}

func (a *App) Run() {
	errChan := make(chan error, 2)

	// Запуск gRPC
	go func() {
		if err := a.gRPCApp.Run(); err != nil {
			errChan <- fmt.Errorf("gRPC failed: %w", err)
		}
	}()

	// Запуск HTTP
	go func() {
		if err := a.httpApp.Run(); err != nil {
			errChan <- fmt.Errorf("HTTP failed: %w", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	select {
	case err := <-errChan:
		a.log.Error().Err(err).Msg("shortener service failed")
		panic(err)
	default:
		a.log.Info().Msg("both servers started successfully")
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-stop:
		a.log.Info().Msg("shutting down servers")
	case err := <-errChan:
		a.log.Error().Err(err).Msg("server failed")
		panic(err)
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	a.gRPCApp.Shutdown()

	if err := a.httpApp.Shutdown(ctx); err != nil {
		a.log.Error().Err(err).Msg("HTTP shutdown error")
	}

	a.log.Info().Msg("shortener service gracefully shutdown")
}
