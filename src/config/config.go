package config

import (
	"log"

	"github.com/go-playground/validator"
	"github.com/spf13/viper"
)

type Config struct {
	App AppConfig `mapstructure:",squash"`
	DB  DBConfig  `mapstructure:",squash"`
}

type AppConfig struct {
	Name    string `mapstructure:"APP_NAME" validate:"required"`
	Version string `mapstructure:"APP_VERSION" validate:"required"`
}

type DBConfig struct {
	Driver   string `mapstructure:"DB_DRIVER" validate:"required"`
	Host     string `mapstructure:"DB_HOST" validate:"required"`
	Port     string `mapstructure:"DB_PORT" validate:"required"`
	User     string `mapstructure:"DB_USER" validate:"required"`
	Password string `mapstructure:"DB_PASSWORD" validate:"required"`
	Name     string `mapstructure:"DB_NAME" validate:"required"`
	SSLMode  string `mapstructure:"DB_SSL_MODE" validate:"required"`
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
