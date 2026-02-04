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

func (s *Service) ProcessWebhook(ctx context.Context, p AvitoWebhookPayload) {
	log.Printf("PROCESSING WEBHOOK: %d messages\n", len(p.Messages))

	for i, msg := range p.Messages {
		log.Printf("MSG #%d TYPE: %s TEXT: %s\n", i, msg.Type, msg.Content.Text)

		if msg.Type == "system" && msg.Content.Text != "" {
			log.Println("→ SENDING TO TELEGRAM")
			if err := s.tg.Send(msg.Content.Text); err != nil {
				log.Println("TG SEND ERROR:", err)
			} else {
				log.Println("TG SEND OK")
			}
		}
	}
}
