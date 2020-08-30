package main

import (
	"context"
	"github.com/gumeniukcom/achecker/configs"
	"github.com/gumeniukcom/achecker/postgres"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg := configs.ReadConfig("./config.toml")

	db, err := postgres.New(cfg.Postgresql)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed to init pg")
	}
	query := "CREATE TABLE IF NOT EXISTS checks (" +
		"id serial, " +
		"domain text, " +
		"status_code int," +
		"error text," +
		"created_on timestamp);" +
		"alter table checks  add constraint checks_pk  primary key (id);"

	_, err = db.Exec(context.Background(), query)
	if err != nil {
		log.Error().Err(err).Msg("error on exec")
	}

}
