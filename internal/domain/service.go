package domain

import (
	"context"
	"fmt"
	"log"
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

const accountAlias = "Самара Jaecoo"

func (s *Service) ProcessWebhook(ctx context.Context, evt AvitoWebhook) {
	v := evt.Payload.Value

	text := v.Content.Text
	if text == "" {
		log.Println("SKIP: empty text")
		return
	}

	out := fmt.Sprintf("Аккаунт %s:\n%s", accountAlias, text)

	// мягкая пометка, НИЧЕГО не фильтруем
	if v.FlowID != "" {
		out += "\n\n(вероятно системное)"
	}

	log.Println("→ SENDING TO TELEGRAM")
	if err := s.tg.Send(out); err != nil {
		log.Println("TG SEND ERROR:", err)
	} else {
		log.Println("TG SEND OK")
	}
}
