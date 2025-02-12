package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	AppVersion       string `mapstructure:"APP_VERSION"`
	AppPort          string `mapstructure:"APP_PORT"`
	AppEnv           string `mapstructure:"APP_ENV"`
	AppName          string `mapstructure:"APP_NAME"`
	AppHost          string `mapstructure:"APP_HOST"`
	AppEnablePrefork bool   `mapstructure:"APP_ENABLE_PREFORK"`
	AppCORSWhitelist string `mapstructure:"APP_CORS_WHITELIST"`
	AppTimezone      *time.Location

	DatabasePostgresHost     string `mapstructure:"DATABASE_POSTGRES_HOST"`
	DatabasePostgresPort     string `mapstructure:"DATABASE_POSTGRES_PORT"`
	DatabasePostgresName     string `mapstructure:"DATABASE_POSTGRES_NAME"`
	DatabasePostgresUsername string `mapstructure:"DATABASE_POSTGRES_USERNAME"`
	DatabasePostgresPassword string `mapstructure:"DATABASE_POSTGRES_PASSWORD"`

	HXMSConfig
}

type HXMSConfig struct {
	HXMSHost string `mapstructure:"HXMS_HOST"`
}

func LoadConfig(path, env string) *Config {
	fmt.Println("Loading cfg with viper...")
	viperInstance := viper.New()
	viperInstance.SetConfigFile(fmt.Sprintf("%sconfig.%s.env", path, env))
	viperInstance.AutomaticEnv()

	if err := viperInstance.ReadInConfig(); err != nil {
		panic(err)
	}

	cfg := new(Config)
	if err := viperInstance.Unmarshal(cfg); err != nil {
		panic(err)
	}
	cfg.AppEnv = env

	timezone, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		panic(err)
	}

	cfg.AppTimezone = timezone
	return cfg
}
