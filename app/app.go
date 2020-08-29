package app

import (
	"github.com/gumeniukcom/achecker/configs"
	"github.com/gumeniukcom/achecker/signals"

	"github.com/rs/zerolog/log"
)

// App container for app
type App struct {
	cfg *configs.Config
}

// New return new instance of App
func New(cfg *configs.Config) *App {
	return &App{
		cfg: cfg,
	}
}

// Run run application
func (app *App) Run() error {

	select {
	case <-signals.WaitExit():
		return app.Stop()
	}

	return nil
}

// Stop stop application
func (app *App) Stop() error {
	log.Info().
		Msg("trying to stop application")
	return nil
}
