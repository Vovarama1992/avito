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

func (s *Service) ProcessWebhook(ctx context.Context, evt AvitoWebhook) {
	// фильтр в сервисе: служебное = author_id == user_id
	v := evt.Payload.Value
	if v.AuthorID != v.UserID {
		log.Println("SKIP: author_id != user_id")
		return
	}

	text := v.Content.Text
	if text == "" {
		log.Println("SKIP: empty text")
		return
	}

	out := fmt.Sprintf("Аккаунт %d: %s", v.UserID, text)

	log.Println("→ SENDING TO TELEGRAM")
	if err := s.tg.Send(out); err != nil {
		log.Println("TG SEND ERROR:", err)
	} else {
		log.Println("TG SEND OK")
	}
}
