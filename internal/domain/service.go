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
	log.Println("→ SENDING RAW TO TELEGRAM")
	_ = s.tg.Send(text)
}
