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
	alias  string
}

func NewService(sender Sender, alias string) *Service {
	return &Service{
		sender: sender,
		alias:  alias,
	}
}

func (s *Service) ProcessWebhook(ctx context.Context, evt AvitoWebhook) {
	v := evt.Payload.Value

	text := v.Content.Text
	if text == "" {
		log.Println("SKIP: empty text")
		return
	}

	out := fmt.Sprintf("Аккаунт %s:\n%s", s.alias, text)

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

func (s *Service) ProcessSystemMessage(ctx context.Context, text string) {
	if text == "" {
		log.Println("SKIP: empty system text")
		return
	}

	out := fmt.Sprintf("Аккаунт %s:\n%s", s.alias, text)

	log.Println("→ SENDING POLLED SYSTEM MESSAGE TO MATRIX")
	if err := s.sender.Send(out); err != nil {
		log.Println("MATRIX SEND ERROR:", err)
	} else {
		log.Println("MATRIX SEND OK")
	}
}
