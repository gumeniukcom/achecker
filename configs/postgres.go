package configs

// PostgresqlConf postgresql config
type PostgresqlConf struct {
	ConnectString string `mapstructure:"CONNECT_STRING" toml:"connect_string"`
}
