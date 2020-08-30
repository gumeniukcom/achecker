package configs

// CheckerConf contains a basic checker conf
type CheckerConf struct {
	TimeoutSecond int  `mapstructure:"TIMEOUT" toml:"timeout"`
	Normalize     bool `mapstructure:"NORMALIZE" toml:"normalize"`
}
