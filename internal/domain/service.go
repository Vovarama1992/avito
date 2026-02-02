package domain

import (
	"context"
)

type TelegramSender interface {
	Send(text string) error
}

type Service struct {
	tg TelegramSender
}

func NewService(tg TelegramSender) *Service {
	return &Service{tg: tg}
}

func (s *Service) ProcessWebhook(ctx context.Context, p AvitoWebhookPayload) {
	for _, msg := range p.Messages {
		if msg.Type == "system" && msg.Content.Text != "" {
			_ = s.tg.Send(msg.Content.Text)
		}
	}
}
