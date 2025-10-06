package handlers

import (
	"fmt"
	"telegram-shop-bot/internal/services"
	"telegram-shop-bot/pkg/keyboards"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type AdminHandler struct {
	api            *tgbotapi.BotAPI
	productService *services.ProductService
	orderService   *services.OrderService
}

func NewAdminHandler(api *tgbotapi.BotAPI, ps *services.ProductService, os *services.OrderService) *AdminHandler {
	return &AdminHandler{
		api:            api,
		productService: ps,
		orderService:   os,
	}
}

func (h *AdminHandler) HandleAdminCommand(update tgbotapi.Update) bool {
	text := update.Message.Text
	chatID := update.Message.Chat.ID

	switch text {
	case "/admin":
		h.showAdminPanel(chatID)
		return true
	case "📋 Все товары":
		h.showAllProducts(chatID)
		return true
	case "🔙 Вернуться в магазин":
		msg := tgbotapi.NewMessage(chatID, "👋 Главное меню")
		msg.ReplyMarkup = keyboards.MainMenu()
		h.api.Send(msg)
		return true
	}
	return false
}

func (h *AdminHandler) showAdminPanel(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "🔧 *Панель администратора*")
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboards.AdminMenu()
	h.api.Send(msg)
}

func (h *AdminHandler) showAllProducts(chatID int64) {
	products, err := h.productService.GetAll()
	if err != nil || len(products) == 0 {
		msg := tgbotapi.NewMessage(chatID, "📋 Товары отсутствуют")
		h.api.Send(msg)
		return
	}

	text := "📋 *Все товары:*\n\n"
	for _, p := range products {
		text += fmt.Sprintf("ID: %d\n*%s*\n💰 %.2f сом | 📦 %d шт.\n\n",
			p.ID, p.Name, p.Price, p.Stock)
	}

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	h.api.Send(msg)
}
