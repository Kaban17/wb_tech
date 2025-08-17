package config

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Config stores the configuration for the cut utility.
type Config struct {
	Fields    string `mapstructure:"fields"`
	Delimiter string `mapstructure:"delimiter"`
	Separated bool   `mapstructure:"separated"`
}

// New creates a new Config object from command-line flags.
func New() *Config {
	pflag.StringP("fields", "f", "", "select only these fields")
	pflag.StringP("delimiter", "d", "\t", "use DELIM instead of TAB for field delimiter")
	pflag.BoolP("separated", "s", false, "do not print lines not containing delimiters")

	pflag.Parse()

	v := viper.New()
	if err := v.BindPFlags(pflag.CommandLine); err != nil {
		fmt.Fprintf(os.Stderr, "Error binding flags: %v\n", err)
		os.Exit(1)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error unmarshaling config: %v\n", err)
		os.Exit(1)
	}

	return &cfg
}
