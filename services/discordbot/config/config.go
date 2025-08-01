package config

import (
	"log"
	"os"
)

// Config holds all configuration for the application.
type Config struct {
	BotToken string
	GuildID  string
}

// Load reads configuration from environment variables.
// It no longer reads from a .env file directly; this is expected to be handled
// by the execution environment (e.g., Docker Compose).
func Load() (*Config, error) {
	cfg := &Config{
		BotToken: os.Getenv("BOT_TOKEN"),
		GuildID:  os.Getenv("GUILD_ID"),
	}

	if cfg.BotToken == "" {
		log.Fatal("Error: BOT_TOKEN environment variable is not set.")
	}
	if cfg.GuildID == "" {
		log.Fatal("Error: GUILD_ID environment variable is not set.")
	}

	return cfg, nil
} 