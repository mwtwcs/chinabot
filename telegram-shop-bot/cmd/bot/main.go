package main

import (
	"log"
	"telegram-shop-bot/internal/bot"
	"telegram-shop-bot/internal/config"
	"telegram-shop-bot/internal/database"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Инициализация базы данных
	db, err := database.InitDB(cfg.DatabasePath)
	if err != nil {
		log.Fatal("Failed to init database:", err)
	}
	defer db.Close()

	// Запуск миграций
	if err := database.RunMigrations(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Создание и запуск бота
	shopBot, err := bot.NewBot(cfg, db)
	if err != nil {
		log.Fatal("Failed to create bot:", err)
	}

	log.Printf("Bot started: @%s", shopBot.GetUsername())
	shopBot.Start()
}
