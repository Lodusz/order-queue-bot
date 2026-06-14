package telegram

import (
	"context"
	"fmt"
	"order-queue-bot/internal/domain/order/service"
	"strings"

	"gopkg.in/telebot.v3"
)

type Handler struct {
	svc *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(bot *telebot.Bot) {
	bot.Handle("/create", h.HandleCreate)
	bot.Handle("/queue", h.HandleQueue)
}

func (h *Handler) HandleCreate(c telebot.Context) error {
	description := strings.Join(c.Args(), " ")
	ord, err := h.svc.CreateOrder(context.Background(), c.Sender().ID, description)
	if err != nil {
		return c.Send("Ошибка: " + err.Error())
	}
	return c.Send(fmt.Sprintf("✅ Заказ #%d создан!", ord.ID))
}

func (h *Handler) HandleQueue(c telebot.Context) error {
	orders, _ := h.svc.GetQueue(context.Background())
	var sb strings.Builder
	for _, o := range orders {
		sb.WriteString(fmt.Sprintf("ID: %d, Описание: %s\n", o.ID, o.Description))
	}
	return c.Send(sb.String())
}
