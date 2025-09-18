package app

import (
	"context"
	"github.com/nikitamavrenko/telegram-bot-service/internal/config"
	"github.com/nikitamavrenko/telegram-bot-service/internal/grpc/shortener"
	"github.com/nikitamavrenko/telegram-bot-service/internal/telegram"
	"github.com/rs/zerolog"
)

type App struct {
	log             *zerolog.Logger
	shortenerClient *shortener.Client
	tgbot           *telegram.Bot
}

func New(cfg *config.Config, log *zerolog.Logger) (*App, error) {
	shortenerContext := context.Background()
	shortenerClient, err := shortener.NewClient(shortenerContext, cfg.ShortenerClient.Addr)
	if err != nil {
		return nil, err
	}

	tgbot, err := telegram.NewBot(cfg.Tg.Token, shortenerClient)
	if err != nil {
		return nil, err
	}

	return &App{
		log:             log,
		shortenerClient: shortenerClient,
		tgbot:           tgbot,
	}, nil
}

func (a *App) MustRun() {
	if err := a.tgbot.Start(); err != nil {
		panic(err)
	}
}
