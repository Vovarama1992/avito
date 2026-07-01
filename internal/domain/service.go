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
	sender                Sender
	alias                 string
	webhookForwardEnabled bool
}

func NewService(sender Sender, alias string, webhookForwardEnabled bool) *Service {
	return &Service{
		sender:                sender,
		alias:                 alias,
		webhookForwardEnabled: webhookForwardEnabled,
	}
}

func (s *Service) ProcessWebhook(ctx context.Context, evt AvitoWebhook) {
	v := evt.Payload.Value

	if !s.webhookForwardEnabled {
		log.Printf(
			"WEBHOOK RECEIVED BUT NOT SENT: reason=webhook_forward_disabled id=%s chat_id=%s msg_type=%s text=%q",
			v.ID,
			v.ChatID,
			v.Type,
			preview(v.Content.Text),
		)
		return
	}

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

func preview(text string) string {
	const limit = 160
	if len(text) <= limit {
		return text
	}
	return text[:limit] + "..."
}

func (s *Service) ProcessSystemMessage(ctx context.Context, source, accountID, chatID, text string) {
	if text == "" {
		log.Println("SKIP: empty system text")
		return
	}

	out := fmt.Sprintf(
		"Аккаунт %s:\nИсточник: %s\naccount_id: %s\nchat_id: %s\n%s",
		s.alias,
		source,
		accountID,
		chatID,
		text,
	)

	log.Println("→ SENDING POLLED SYSTEM MESSAGE TO MATRIX")
	if err := s.sender.Send(out); err != nil {
		log.Println("MATRIX SEND ERROR:", err)
	} else {
		log.Println("MATRIX SEND OK")
	}
}
