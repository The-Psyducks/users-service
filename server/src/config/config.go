package config

import (
	"github.com/spf13/viper"
	"log/slog"
)

// Config struct that holds the configuration of the server
type Config struct {
	Host             string
	Port             string
	Environment      string
	DatabaseHost     string
	DatabasePort     string
	DatabaseName     string
	DatabaseUser     string
	DatabasePassword string
}

// LoadConfig loads the configuration from the Environment variables
func LoadConfig() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	viper.SetDefault("PORT", "8080")
	viper.SetDefault("HOST", "0.0.0.0")
	viper.SetDefault("ENVIRONMENT", "development")

	if err := viper.ReadInConfig(); err != nil {
		slog.Error(".env configuration error", slog.String("error", err.Error()), slog.String("info:", "The .env configuration file was not found or there was an error reading it"))
	}

	return &Config{
		Host:             viper.GetString("HOST"),
		Port:             viper.GetString("PORT"),
		Environment:      viper.GetString("ENVIRONMENT"),
		DatabaseHost:     viper.GetString("DATABASE_HOST"),
		DatabasePort:     viper.GetString("DATABASE_PORT"),
		DatabaseName:     viper.GetString("DATABASE_NAME"),
		DatabaseUser:     viper.GetString("DATABASE_USER"),
		DatabasePassword: viper.GetString("DATABASE_PASSWORD"),
	}
}