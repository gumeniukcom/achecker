package configs

type InitialOffset string

const (
	// OffsetNewest stands for the log head offset, i.e. the offset that will be
	// assigned to the next message that will be produced to the partition.
	OffsetNewest = InitialOffset("newest")
	// OffsetOldest stands for the oldest offset available on the broker for a
	// partition.
	OffsetOldest = InitialOffset("oldest")
)

// KafkaConf ...
type KafkaConf struct {
	Brokers       []string      `mapstructure:"BROKERS" toml:"brokers"`
	Group         string        `mapstructure:"GROUP" toml:"group"`
	Debug         bool          `mapstructure:"DEBUG" toml:"debug"`
	InitialOffset InitialOffset `mapstructure:"INITIAL_OFFSET" toml:"initial_offset"`
	SSl           bool          `mapstructure:"SSL" toml:"ssl"`
	FileCAPath    string        `mapstructure:"CAPATH" toml:"capath"`
	FileCertPath  string        `mapstructure:"CERTPATH" toml:"certpath"`
	FileKeyPath   string        `mapstructure:"KEYPATH" toml:"keypath"`
	Version       string        `mapstructure:"VERSION" toml:"version"`
}
