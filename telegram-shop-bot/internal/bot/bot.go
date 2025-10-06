package bot

import (
	"database/sql"
	"telegram-shop-bot/internal/config"
	"telegram-shop-bot/internal/handlers"
	"telegram-shop-bot/internal/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api             *tgbotapi.BotAPI
	config          *config.Config
	userHandler     *handlers.UserHandler
	adminHandler    *handlers.AdminHandler
	callbackHandler *handlers.CallbackHandler
	adminIDs        map[int64]bool
}

func NewBot(cfg *config.Config, db *sql.DB) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		return nil, err
	}

	api.Debug = cfg.Debug

	// Инициализация сервисов
	productService := services.NewProductService(db)
	cartService := services.NewCartService(db)
	orderService := services.NewOrderService(db)

	// Создание админ ID map
	adminIDs := make(map[int64]bool)
	for _, id := range cfg.AdminIDs {
		adminIDs[id] = true
	}

	bot := &Bot{
		api:      api,
		config:   cfg,
		adminIDs: adminIDs,
	}

	// Инициализация handlers
	bot.userHandler = handlers.NewUserHandler(api, productService, cartService, orderService)
	bot.adminHandler = handlers.NewAdminHandler(api, productService, orderService)
	bot.callbackHandler = handlers.NewCallbackHandler(api, cartService, orderService, productService)

	return bot, nil
}

func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			b.handleMessage(update)
		} else if update.CallbackQuery != nil {
			b.callbackHandler.Handle(update)
		}
	}
}

func (b *Bot) handleMessage(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	isAdmin := b.adminIDs[chatID]

	if isAdmin && b.adminHandler.HandleAdminCommand(update) {
		return
	}

	b.userHandler.Handle(update)
}

func (b *Bot) GetUsername() string {
	return b.api.Self.UserName
}
