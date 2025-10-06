package handlers

import (
	"fmt"
	"strconv"
	"strings"
	"telegram-shop-bot/internal/services"
	"telegram-shop-bot/pkg/keyboards"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CallbackHandler struct {
	api            *tgbotapi.BotAPI
	cartService    *services.CartService
	orderService   *services.OrderService
	productService *services.ProductService
}

func NewCallbackHandler(api *tgbotapi.BotAPI, cs *services.CartService, os *services.OrderService, ps *services.ProductService) *CallbackHandler {
	return &CallbackHandler{
		api:            api,
		cartService:    cs,
		orderService:   os,
		productService: ps,
	}
}

func (h *CallbackHandler) Handle(update tgbotapi.Update) {
	callback := update.CallbackQuery
	data := callback.Data
	chatID := callback.Message.Chat.ID

	if strings.HasPrefix(data, "buy_") {
		h.handleBuy(callback, chatID, data)
	} else if data == "checkout" {
		h.handleCheckout(callback, chatID)
	} else if data == "clear_cart" {
		h.handleClearCart(callback, chatID)
	} else if strings.HasPrefix(data, "pay_") {
		h.handlePayment(callback, chatID, data)
	}
}

func (h *CallbackHandler) handleBuy(callback *tgbotapi.CallbackQuery, chatID int64, data string) {
	productID, _ := strconv.Atoi(strings.TrimPrefix(data, "buy_"))
	h.cartService.AddItem(chatID, productID, 1)

	answer := tgbotapi.NewCallback(callback.ID, "✅ Добавлено в корзину!")
	h.api.Request(answer)
}

func (h *CallbackHandler) handleCheckout(callback *tgbotapi.CallbackQuery, chatID int64) {
	items, _ := h.cartService.GetItems(chatID)
	var total float64

	for _, item := range items {
		product, _ := h.productService.GetByID(item.ProductID)
		if product != nil {
			orderTotal := product.Price * float64(item.Quantity)
			total += orderTotal
			h.orderService.Create(chatID, product.ID, item.Quantity, orderTotal)
		}
	}

	h.cartService.Clear(chatID)

	text := fmt.Sprintf("✅ *Заказ оформлен!*\n\n💰 Сумма: %.2f сом\n\nВыберите способ оплаты:", total)
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboards.PaymentMethods()
	h.api.Send(msg)
}

func (h *CallbackHandler) handleClearCart(callback *tgbotapi.CallbackQuery, chatID int64) {
	h.cartService.Clear(chatID)

	answer := tgbotapi.NewCallback(callback.ID, "🗑 Корзина очищена")
	h.api.Request(answer)

	msg := tgbotapi.NewMessage(chatID, "🛒 Корзина пуста")
	h.api.Send(msg)
}

func (h *CallbackHandler) handlePayment(callback *tgbotapi.CallbackQuery, chatID int64, data string) {
	payMethod := "картой"
	if data == "pay_cash" {
		payMethod = "наличными при получении"
	}

	text := fmt.Sprintf("✅ Заказ принят!\n\nОплата: %s\n\nМы свяжемся с вами в ближайшее время!", payMethod)
	msg := tgbotapi.NewMessage(chatID, text)
	h.api.Send(msg)

	answer := tgbotapi.NewCallback(callback.ID, "✅")
	h.api.Request(answer)
}
