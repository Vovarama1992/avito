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
	log.Println("PROCESS MESSAGE TEXT:", text)

	if text == "" {
		log.Println("EMPTY TEXT — SKIP")
		return
	}

	log.Println("→ SENDING TO TELEGRAM")

	if err := s.tg.Send(text); err != nil {
		log.Println("TG SEND ERROR:", err)
	} else {
		log.Println("TG SEND OK")
	}
}
