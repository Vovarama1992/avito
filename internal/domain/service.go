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

// хардкод маппинга
var accountAlias = map[int64]string{
	375283938: "Самара Jaecoo",
}

func (s *Service) ProcessWebhook(ctx context.Context, evt AvitoWebhook) {
	v := evt.Payload.Value

	// служебное
	if v.AuthorID != v.UserID {
		log.Println("SKIP: author_id != user_id")
		return
	}

	text := v.Content.Text
	if text == "" {
		log.Println("SKIP: empty text")
		return
	}

	alias, ok := accountAlias[v.UserID]
	if !ok {
		alias = fmt.Sprintf("%d", v.UserID)
	}

	out := fmt.Sprintf(
		"Аккаунт %s:\n%s",
		alias,
		text,
	)

	log.Println("→ SENDING TO TELEGRAM")
	if err := s.tg.Send(out); err != nil {
		log.Println("TG SEND ERROR:", err)
	} else {
		log.Println("TG SEND OK")
	}
}
