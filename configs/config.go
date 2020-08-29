package configs

import (
	"os"

	"github.com/jacobstr/confer"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Config ...
type Config struct {
	Env     string     `mapstructure:"ENV" toml:"env"`
	AppName string     `mapstructure:"APPNAME" toml:"appName"`
	Version string     `mapstructure:"VERSION" toml:"version"`
	Logger  LoggerConf `mapstructure:"LOGGER" toml:"logger"`
}

// ReadConfig tryes to load config from Env, if can't then Toml
func ReadConfig(filename string) *Config {
	c := confer.NewConfig()
	if err := c.ReadPaths(filename); err != nil {
		log.Fatal().
			Err(err).
			Str("config_file_name", filename).
			Msg("can not read config file")
	}

	setupDefaultConfigParameters(c)
	c.AutomaticEnv()

	cfg, err := createConfigFromConfer(c)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("can not create config")
	}

	zerolog.SetGlobalLevel(cfg.Logger.Level)
	if cfg.Logger.OutputType == LoggerConsoleOutput {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	}
	return &cfg
}

func setupDefaultConfigParameters(c *confer.Config) {
	c.SetDefault("env", "dev")

	// Logger block
	c.SetDefault("logger.level", 0)
	c.SetDefault("logger.output_type", LoggerJsonOutput)
	c.SetDefault("logger.time_field_format", "")

}

func createConfigFromConfer(c *confer.Config) (Config, error) {
	cfg := Config{}

	// Commons
	cfg.Env = c.GetString("env")
	cfg.AppName = c.GetString("appName")
	cfg.Version = c.GetString("version")

	// Logger block
	cfg.Logger.Level = zerolog.Level(c.GetInt("logger.level"))
	cfg.Logger.OutputType = OutputType(c.GetString("logger.output_type"))
	cfg.Logger.TimeFieldFormat = c.GetString("logger.time_field_format")

	return cfg, nil
}
