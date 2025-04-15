package config

import (
	"log"

	"github.com/go-playground/validator"
	"github.com/spf13/viper"
)

type Config struct {
	App AppConfig `mapstructure:",squash"`
}

type AppConfig struct {
	Name    string `mapstructure:"APP_NAME" validate:"required"`
	Version string `mapstructure:"APP_VERSION" validate:"required"`
}

func (c *Config) validate() error {
	return validator.New().Struct(c)
}

func Load() *Config {
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		log.Fatalf("Error unmarshalling config: %v", err)
	}

	if err := config.validate(); err != nil {
		log.Fatalf("Error validating config: %v", err)
	}

	return config
}
