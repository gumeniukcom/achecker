package app

import (
	"github.com/gumeniukcom/achecker/checkdaemon"
	"github.com/gumeniukcom/achecker/configs"
	"github.com/gumeniukcom/achecker/signals"
	"github.com/rs/zerolog/log"
)

// App container for app
type App struct {
	cfg         *configs.Config
	chechdaemon *checkdaemon.Daemon
}

// New return new instance of App
func New(cfg *configs.Config) *App {

	return &App{
		cfg:         cfg,
		chechdaemon: checkdaemon.New(cfg),
	}
}

// Run run application
func (app *App) Run() error {

	if err := app.chechdaemon.Run(); err != nil {
		log.Error().
			Err(err).
			Msg("error on check checkdaemon")
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
	app.chechdaemon.Stop()
}
