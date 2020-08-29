package configs

import "github.com/rs/zerolog"

type OutputType string

const (
	LoggerJsonOutput    = OutputType("json")
	LoggerConsoleOutput = OutputType("console")
)

// LoggerConf contains a basic logger parameters
type LoggerConf struct {
	Level           zerolog.Level `mapstructure:"LEVEL" toml:"level"`
	OutputType      OutputType    `mapstructure:"OUTPUT_TYPE" tomal:"output_type"`
	TimeFieldFormat string        `mapstructure:"TIME_FIELD_FORMAT" toml:"time_field_format"`
}
