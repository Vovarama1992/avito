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

func (s *Service) ProcessMessage(ctx context.Context, text string) {
	log.Println("SERVICE GOT TEXT:", text)

	if text == "" {
		return
	}

	if err := s.tg.Send(text); err != nil {
		log.Println("TG SEND ERROR:", err)
	} else {
		log.Println("TG SEND OK")
	}
}
