package configs

// ResultDaemonConf contains a basic resultdaemon conf
type ResultDaemonConf struct {
	ResultTopic string `mapstructure:"RESULTTOPIC" toml:"result_topic"`
	KafkaGroup  string `mapstructure:"KAFKAGROUP" toml:"kafka_group"`
}
