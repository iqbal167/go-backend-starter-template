package config

import (
	"fmt"

	"github.com/go-playground/validator"
	"github.com/spf13/viper"
)

type Config struct {
	App        AppConfig  `mapstructure:",squash"`
	DB         DBConfig   `mapstructure:",squash"`
	CORSConfig CORSConfig `mapstructure:",squash"`
}

type AppConfig struct {
	Name    string `mapstructure:"APP_NAME"`
	Version string `mapstructure:"APP_VERSION"`
	Port    int    `mapstructure:"APP_PORT"`
}

type CORSConfig struct {
	AllowedOrigins   []string `mapstructure:"CORS_ALLOWED_ORIGINS"`
	AllowedMethods   []string `mapstructure:"CORS_ALLOWED_METHODS"`
	AllowedHeaders   []string `mapstructure:"CORS_ALLOWED_HEADERS"`
	AllowCredentials bool     `mapstructure:"CORS_ALLOW_CREDENTIALS"`
	MaxAge           int      `mapstructure:"CORS_MAX_AGE"`
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

func newConfig() (*Config, error) {
	config := &Config{}
	config.setDefaults()
	config.configureViper()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("could not read config file. %w", err)
	}

	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to parse config values. %w", err)
	}

	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("invalid config values. %w", err)
	}

	return config, nil
}

func (c *Config) configureViper() {
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
}

func (c *Config) setDefaults() {
	viper.SetDefault("APP_NAME", "server")
	viper.SetDefault("APP_VERSION", "1.0.0")
	viper.SetDefault("APP_PORT", 8080)
	viper.SetDefault("CORS_ALLOWED_ORIGINS", []string{"*"})
	viper.SetDefault("CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"})
	viper.SetDefault("CORS_ALLOWED_HEADERS", []string{"Accept", "Authorization", "Content-Type"})
	viper.SetDefault("CORS_ALLOW_CREDENTIALS", false)
	viper.SetDefault("CORS_MAX_AGE", 600)

}

func (c *Config) validate() error {
	return validator.New().Struct(c)
}

func Load() (*Config, error) {
	config, err := newConfig()
	if err != nil {
		return nil, err
	}
	return config, nil
}
