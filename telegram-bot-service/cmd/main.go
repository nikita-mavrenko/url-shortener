package main

import (
	"github.com/nikitamavrenko/telegram-bot-service/internal/app"
	"github.com/nikitamavrenko/telegram-bot-service/internal/config"
	"github.com/rs/zerolog"
	"os"
)

func main() {
	log := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()

	log.Info().Msg("start app")
	cfg := config.MustLoad()

	a, err := app.New(cfg, &log)
	if err != nil {
		log.Fatal().Err(err).Msg("app initialization failed")
	}
	a.MustRun()
}
