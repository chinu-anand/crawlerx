package configs

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port        string
	DatabaseURL string
	RedisURL    string
}

var App Config

func LoadConfig() {
	viper.SetConfigFile(".env")

	// Defaults
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("DATABASE_URL", "host=localhost user=postgres password=postgres dbname=crawlerx port=5432 sslmode=disable")
	viper.SetDefault("REDIS_URL", "redis://localhost:6379/0")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Failed to load config: %v", err)
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	App = Config{
		Port:        viper.GetString("PORT"),
		DatabaseURL: viper.GetString("DATABASE_URL"),
		RedisURL:    viper.GetString("REDIS_URL"),
	}
}
