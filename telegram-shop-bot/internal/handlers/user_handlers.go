package handlers

import (
	"fmt"
	"telegram-shop-bot/internal/services"
	"telegram-shop-bot/pkg/keyboards"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UserHandler struct {
	api            *tgbotapi.BotAPI
	productService *services.ProductService
	cartService    *services.CartService
	orderService   *services.OrderService
}

func NewUserHandler(api *tgbotapi.BotAPI, ps *services.ProductService, cs *services.CartService, os *services.OrderService) *UserHandler {
	return &UserHandler{
		api:            api,
		productService: ps,
		cartService:    cs,
		orderService:   os,
	}
}

func (h *UserHandler) Handle(update tgbotapi.Update) {
	text := update.Message.Text
	chatID := update.Message.Chat.ID

	switch text {
	case "/start":
		h.handleStart(chatID)
	case "🛍 Каталог":
		h.showCatalog(chatID)
	case "🛒 Корзина":
		h.showCart(chatID)
	case "📦 Мои заказы":
		h.showOrders(chatID)
	case "ℹ️ О магазине":
		h.showAbout(chatID)
	}
}

func (h *UserHandler) handleStart(chatID int64) {
	text := `👋 *Добро пожаловать в наш магазин!*

🛍 Выбирайте товары
📦 Быстрая доставка
💳 Удобная оплата`

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboards.MainMenu()
	h.api.Send(msg)
}

func (h *UserHandler) showCatalog(chatID int64) {
	products, err := h.productService.GetAll()
	if err != nil || len(products) == 0 {
		msg := tgbotapi.NewMessage(chatID, "❌ Каталог пуст")
		h.api.Send(msg)
		return
	}

	text := "🛍 *Наш каталог:*\n\n"
	var buttons [][]tgbotapi.InlineKeyboardButton

	for _, p := range products {
		text += fmt.Sprintf("*%s*\n%s\n💰 %.2f сом | 📦 %d шт.\n\n",
			p.Name, p.Description, p.Price, p.Stock)

		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("🛒 %s", p.Name),
				fmt.Sprintf("buy_%d", p.ID),
			),
		)
		buttons = append(buttons, row)
	}

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(buttons...)
	h.api.Send(msg)
}

func (h *UserHandler) showCart(chatID int64) {
	items, err := h.cartService.GetItems(chatID)
	if err != nil || len(items) == 0 {
		msg := tgbotapi.NewMessage(chatID, "🛒 Корзина пуста")
		h.api.Send(msg)
		return
	}

	text := "🛒 *Ваша корзина:*\n\n"
	var total float64

	for _, item := range items {
		product, _ := h.productService.GetByID(item.ProductID)
		if product != nil {
			itemTotal := product.Price * float64(item.Quantity)
			total += itemTotal
			text += fmt.Sprintf("• %s x%d = %.2f сом\n", product.Name, item.Quantity, itemTotal)
		}
	}

	text += fmt.Sprintf("\n💰 *Итого: %.2f сом*", total)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Оформить", "checkout"),
			tgbotapi.NewInlineKeyboardButtonData("🗑 Очистить", "clear_cart"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	h.api.Send(msg)
}

func (h *UserHandler) showOrders(chatID int64) {
	orders, err := h.orderService.GetUserOrders(chatID)
	if err != nil || len(orders) == 0 {
		msg := tgbotapi.NewMessage(chatID, "📦 У вас пока нет заказов")
		h.api.Send(msg)
		return
	}

	text := "📦 *Ваши заказы:*\n\n"
	for _, order := range orders {
		product, _ := h.productService.GetByID(order.ProductID)
		if product != nil {
			text += fmt.Sprintf("Заказ #%d\n%s x%d\n💰 %.2f сом\nСтатус: %s\n\n",
				order.ID, product.Name, order.Quantity, order.Total, order.Status)
		}
	}

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	h.api.Send(msg)
}

func (h *UserHandler) showAbout(chatID int64) {
	text := "ℹ️ *О магазине*\n\n🕒 Работаем: 9:00 - 21:00\n📞 Поддержка: @shop_support"
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	h.api.Send(msg)
}
