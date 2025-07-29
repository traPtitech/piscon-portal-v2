package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func LoadFile() (Config, error) {
	viper.SetConfigName("piscon_runner")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("unmarshal config: %w", err)
	}

	return cfg, nil
}

type Portal struct {
	Address string `mapstructure:"address"`
}

type Problem struct {
	Name    string         `mapstructure:"name"`
	Options map[string]any `mapstructure:"options"`
}

type Config struct {
	Portal  Portal  `mapstructure:"portal"`
	Problem Problem `mapstructure:"problem"`
}

func (c Config) Validate() error {
	if c.Portal.Address == "" {
		return fmt.Errorf("portal address is required")
	}
	if c.Problem.Name == "" {
		return fmt.Errorf("problem name is required")
	}

	return nil
}
