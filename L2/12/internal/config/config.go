package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	A           int    `mapstructure:"A"`
	B           int    `mapstructure:"B"`
	C           int    `mapstructure:"C"`
	CountString bool   `mapstructure:"count-string"`
	IgnoreCase  bool   `mapstructure:"ignore-case"`
	Invert      bool   `mapstructure:"invert"`
	FixString   bool   `mapstructure:"fix-string"`
	LineNum     bool   `mapstructure:"line-num"`
	InputFile   string `mapstructure:"input-file"`
	Regex       string `mapstructure:"regex"`
}

func New() *Config {
	v := viper.New()

	v.SetDefault("A", 0)
	v.SetDefault("B", 0)
	v.SetDefault("C", 0)
	v.SetDefault("count-string", false)
	v.SetDefault("ignore-case", false)
	v.SetDefault("invert", false)
	v.SetDefault("fix-string", false)
	v.SetDefault("line-num", false)
	v.SetDefault("input-file", "")
	v.SetDefault("regex", "")

	pflag.IntP("A", "A", v.GetInt("A"), "Print string before match")
	pflag.IntP("B", "B", v.GetInt("B"), "Print string after match")
	pflag.IntP("C", "C", v.GetInt("C"), "Print string before and after match")

	pflag.BoolP("count-string", "c", v.GetBool("count-string"), "Print count of matching lines")
	pflag.BoolP("ignore-case", "i", v.GetBool("ignore-case"), "Ignore case when matching")
	pflag.BoolP("invert", "v", v.GetBool("invert"), "Invert match")
	pflag.BoolP("fix-string", "F", v.GetBool("fix-string"), "Fix string")
	pflag.BoolP("line-num", "n", v.GetBool("line-num"), "Print line number")

	v.BindPFlags(pflag.CommandLine)

	v.SetEnvPrefix("MYAPP")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	pflag.Parse()

	nonFlagArgs := pflag.Args()

	if len(nonFlagArgs) > 0 {
		v.Set("regex", nonFlagArgs[0])
	}
	if len(nonFlagArgs) > 1 {
		v.Set("input-file", nonFlagArgs[len(nonFlagArgs)-1])
	} else if len(nonFlagArgs) == 1 {
	}

	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to decode config into struct: %v\n", err)
		os.Exit(1)
	}

	return cfg
}
