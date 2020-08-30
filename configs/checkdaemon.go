package configs

// CheckDaemonConf contains a basic checkdaemon conf
type CheckDaemonConf struct {
	ResultTopic string `mapstructure:"RESULTTOPIC" toml:"result_topic"`
	CheckTopic  string `mapstructure:"CHECKTOPIC" toml:"check_topic"`
}
