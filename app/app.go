package app

import (
	"github.com/gumeniukcom/achecker/checkdaemon"
	"github.com/gumeniukcom/achecker/configs"
	"github.com/gumeniukcom/achecker/resultdaemon"
	"github.com/gumeniukcom/achecker/signals"
	"github.com/rs/zerolog/log"
)

// App container for app
type App struct {
	cfg          *configs.Config
	checkdaemon  *checkdaemon.Daemon
	resultdaemon *resultdaemon.Resoluter
}

// New return new instance of App
func New(cfg configs.Config) *App {
	return &App{
		cfg:          &cfg,
		checkdaemon:  checkdaemon.New(cfg),
		resultdaemon: resultdaemon.New(cfg),
	}
}

// Run run application
func (app *App) Run() error {

	if err := app.checkdaemon.Run(); err != nil {
		log.Error().
			Err(err).
			Msg("error on check checkdaemon")
		return err
	}

	if err := app.resultdaemon.Run(); err != nil {
		log.Error().
			Err(err).
			Msg("error on check resultdaemon")
		return err
	}

	log.Info().
		Msg("application started")
	select {
	case <-signals.WaitExit():
		app.Stop()
	}

	return nil
}

// Stop stop application
func (app *App) Stop() {
	log.Info().
		Msg("trying to stop application")
	app.checkdaemon.Stop()
	app.resultdaemon.Stop()
}
