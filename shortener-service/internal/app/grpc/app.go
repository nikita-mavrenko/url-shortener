package grpcapp

import (
	"errors"
	"fmt"
	shortenergrpc "github.com/nikitamavrenko/shortener-service/internal/grpc"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"net"
)

type App struct {
	log    *zerolog.Logger
	server *grpc.Server
	port   int
}

func New(log *zerolog.Logger, shortenerService shortenergrpc.Shortener, port int) *App {
	server := grpc.NewServer()

	shortenergrpc.Register(server, shortenerService)

	return &App{
		log:    log,
		server: server,
		port:   port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	a.log.Info().
		Str("addr", l.Addr().String()).
		Int("port", a.port).
		Msg("grpc server starting")

	if err := a.server.Serve(l); err != nil {
		if errors.Is(err, grpc.ErrServerStopped) {
			a.log.Info().Msg("grpc server stopped gracefully")
			return nil
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Shutdown() {
	const op = "grpcapp.Shutdown"

	a.log.Info().Str("op", op).Int("port", a.port).Msg("shutting down grpc server")
	a.server.GracefulStop()
}
