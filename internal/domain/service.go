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

	// системное сообщение от чат-бота Авито: flow_id заполнен
	if v.FlowID == "" {
		log.Println("SKIP: not system (empty flow_id)")
		return
	}

	text := v.Content.Text
	if text == "" {
		log.Println("SKIP: empty text")
		return
	}

	out := fmt.Sprintf("Аккаунт %s:\n%s", accountAlias, text)

	log.Println("→ SENDING TO TELEGRAM")
	if err := s.tg.Send(out); err != nil {
		log.Println("TG SEND ERROR:", err)
	} else {
		log.Println("TG SEND OK")
	}
}
