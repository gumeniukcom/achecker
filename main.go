package main

import (
	"github.com/gumeniukcom/achecker/app"
	"github.com/gumeniukcom/achecker/configs"
	"github.com/rs/zerolog/log"
)

func main() {

	cfg := configs.ReadConfig("./config.toml")

	a := app.New(cfg)
	log.Info().
		Msg("application running")
	if err := a.Run(); err != nil {
		log.Error().Err(err).Msg("some error on run app")
	}

	log.Info().
		Msg("application was stopped")
}
