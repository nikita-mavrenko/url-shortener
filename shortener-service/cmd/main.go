package main

import (
	"github.com/nikitamavrenko/shortener-service/internal/app"
	"github.com/nikitamavrenko/shortener-service/internal/config"
	"github.com/rs/zerolog"
	"os"
)

func main() {

	log := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
	cfg := config.MustLoad()
	log.Info().Msg(cfg.Db.Url)
	application := app.New(&log, cfg)
	application.Run()
}
