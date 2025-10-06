package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken     string
	DatabasePath string
	AdminIDs     []int64
	Debug        bool
}

func LoadConfig() (*Config, error) {
	// Загрузка .env файла
	godotenv.Load()

	adminIDsStr := os.Getenv("ADMIN_IDS")
	adminIDs := parseAdminIDs(adminIDsStr)

	debug, _ := strconv.ParseBool(os.Getenv("DEBUG"))

	return &Config{
		BotToken:     os.Getenv("BOT_TOKEN"),
		DatabasePath: getEnvOrDefault("DATABASE_PATH", "./shop.db"),
		AdminIDs:     adminIDs,
		Debug:        debug,
	}, nil
}

func parseAdminIDs(s string) []int64 {
	var ids []int64
	for _, idStr := range strings.Split(s, ",") {
		if id, err := strconv.ParseInt(strings.TrimSpace(idStr), 10, 64); err == nil {
			ids = append(ids, id)
		}
	}
	return ids
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
