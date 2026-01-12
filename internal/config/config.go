package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

func Init() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("internal/config")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetDefault("server.addr", ":5000")
	viper.SetDefault("db.url", "postgres://postgres:password@localhost:5432/edugov?sslmode=disable")
	viper.SetDefault("jwt.secret", "change-me-in-prod")
	viper.SetDefault("jwt.access_ttl", 2*time.Hour)
	viper.SetDefault("jwt.refresh_ttl", 7*24*time.Hour)

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("read config: %w", err)
	}

	return nil
}
