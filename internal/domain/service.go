package domain

import (
	"context"
	"fmt"
	"log"
)

type Sender interface {
	Send(text string) error
}

type Service struct {
	sender Sender
}

func NewService(sender Sender) *Service {
	return &Service{sender: sender}
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

	if v.FlowID != "" {
		out += "\n\n(вероятно системное)"
	}

	log.Println("→ SENDING TO MATRIX")
	if err := s.sender.Send(out); err != nil {
		log.Println("MATRIX SEND ERROR:", err)
	} else {
		log.Println("MATRIX SEND OK")
	}
}
