package domain

import (
	"context"
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

func (s *Service) ProcessEvent(ctx context.Context, ev AvitoEvent) {
	if ev.Payload.Type != "message" {
		return
	}

	v := ev.Payload.Value

	// служебное: author_id == user_id
	if v.AuthorID != v.UserID {
		return
	}

	if v.Content.Text == "" {
		return
	}

	log.Println("→ SERVICE MESSAGE TO TG:", v.Content.Text)
	_ = s.tg.Send(v.Content.Text)
}
