package config

import (
	"fmt"
	"os"
)

type Config struct {
	BotToken string
	DSN      string
}

func Load() (*Config, error) {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	dsn := os.Getenv("DATABASE_URL")
	if token == "" || dsn == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN or DATABASE_URL not set")
	}
	return &Config{BotToken: token, DSN: dsn}, nil
}
