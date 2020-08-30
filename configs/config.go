package configs

import (
	"os"

	"github.com/jacobstr/confer"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Config ...
type Config struct {
	Env          string           `mapstructure:"ENV" toml:"env"`
	AppName      string           `mapstructure:"APPNAME" toml:"appName"`
	Version      string           `mapstructure:"VERSION" toml:"version"`
	Logger       LoggerConf       `mapstructure:"LOGGER" toml:"logger"`
	Kafka        KafkaConf        `mapstructure:"KAFKA" toml:"kafka"`
	ResultKafka  KafkaConf        `mapstructure:"RESULTKAFKA" toml:"resultkafka"`
	CheckDaemon  CheckDaemonConf  `mapstructure:"CHECKDAEMON" toml:"checkdaemon"`
	Checker      CheckerConf      `mapstructure:"CHECKER" toml:"checker"`
	ResultDaemon ResultDaemonConf `mapstructure:"RESULTDAEMON" toml:"resultdaemon"`
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

	// Kafka block
	c.SetDefault("kafka.initial_offset", OffsetNewest)
	c.SetDefault("kafka.version", "1.0.0.0")

	// ResultKafka block
	c.SetDefault("resultkafka.initial_offset", OffsetNewest)
	c.SetDefault("resultkafka.version", "1.0.0.0")

	// Checkdaemon block

	c.SetDefault("checkdaemon.result_topic", "result")
	c.SetDefault("checkdaemon.check_topic", "check")

	// Checker conf
	c.SetDefault("checker.timeout", 30)
	c.SetDefault("checker.normalize", true)

	// Resultdaemin block

	c.SetDefault("resultdaemon.result_topic", "result")
	c.SetDefault("resultdaemon.kafka_group", "acheckerresult")

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

	// Kafka block
	cfg.Kafka.Brokers = c.GetStringSlice("kafka.brokers")
	cfg.Kafka.Group = c.GetString("kafka.group")
	cfg.Kafka.Debug = c.GetBool("kafka.debug")
	cfg.Kafka.InitialOffset = InitialOffset(c.GetString("kafka.initial_offset"))
	cfg.Kafka.SSl = c.GetBool("kafka.ssl")
	cfg.Kafka.FileCAPath = c.GetString("kafka.capath")
	cfg.Kafka.FileKeyPath = c.GetString("kafka.keypath")
	cfg.Kafka.FileCertPath = c.GetString("kafka.certpath")
	cfg.Kafka.Version = c.GetString("kafka.version")

	// ResultKafka block
	cfg.ResultKafka.Brokers = c.GetStringSlice("resultkafka.brokers")
	cfg.ResultKafka.Group = c.GetString("resultkafka.group")
	cfg.ResultKafka.Debug = c.GetBool("resultkafka.debug")
	cfg.ResultKafka.InitialOffset = InitialOffset(c.GetString("resultkafka.initial_offset"))
	cfg.ResultKafka.SSl = c.GetBool("resultkafka.ssl")
	cfg.ResultKafka.FileCAPath = c.GetString("resultkafka.capath")
	cfg.ResultKafka.FileKeyPath = c.GetString("resultkafka.keypath")
	cfg.ResultKafka.FileCertPath = c.GetString("resultkafka.certpath")
	cfg.ResultKafka.Version = c.GetString("resultkafka.version")

	// Checkdaemon block

	cfg.CheckDaemon.CheckTopic = c.GetString("checkdaemon.check_topic")
	cfg.CheckDaemon.ResultTopic = c.GetString("checkdaemon.result_topic")

	// Checker block

	cfg.Checker.Normalize = c.GetBool("checker.normalize")
	cfg.Checker.TimeoutSecond = c.GetInt("checker.timeout")

	// Resultdaemon block

	cfg.ResultDaemon.ResultTopic = c.GetString("resultdaemon.result_topic")
	cfg.ResultDaemon.KafkaGroup = c.GetString("resultdaemon.kafka_group")

	return cfg, nil
}
