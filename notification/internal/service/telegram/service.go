package telegram

import (
	"bytes"
	"context"
	"embed"
	"text/template"

	"go.uber.org/zap"

	"github.com/voronovsg/rocket-factory/notification/internal/client"
	"github.com/voronovsg/rocket-factory/notification/internal/model"
	"github.com/voronovsg/rocket-factory/platform/pkg/logger"
)

const chatID = 1443445013

//go:embed templates/order_paid_notification.tmpl templates/order_assembled_notification.tmpl
var templateFS embed.FS

var (
	orderPaidData      = template.Must(template.ParseFS(templateFS, "templates/order_paid_notification.tmpl"))
	orderAssembledData = template.Must(template.ParseFS(templateFS, "templates/order_assembled_notification.tmpl"))
)

type service struct {
	telegramClient http.TelegramClient
}

func NewService(telegramClient http.TelegramClient) *service {
	return &service{
		telegramClient: telegramClient,
	}
}

func (s *service) SendOrderPaidNotification(ctx context.Context, event model.OrderPaidEvent) error {
	message, err := s.buildMessage(orderPaidData, event)
	if err != nil {
		return err
	}

	err = s.telegramClient.SendMessage(ctx, chatID, message)
	if err != nil {
		return err
	}

	logger.Info(ctx, "Telegram message sent to chat", zap.Int("chat_id", chatID), zap.String("message", message))
	return nil
}

func (s *service) SendOrderAssembledNotification(ctx context.Context, event model.OrderAssembledEvent) error {
	message, err := s.buildMessage(orderAssembledData, event)
	if err != nil {
		return err
	}

	err = s.telegramClient.SendMessage(ctx, chatID, message)
	if err != nil {
		return err
	}

	logger.Info(ctx, "Telegram message sent to chat", zap.Int("chat_id", chatID), zap.String("message", message))
	return nil
}

func (s *service) buildMessage(tmpl *template.Template, data interface{}) (string, error) {
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
