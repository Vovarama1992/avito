package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Sender struct {
	api     *tgbotapi.BotAPI
	chatIDs []int64
}

func NewSender(api *tgbotapi.BotAPI, chatIDs []int64) *Sender {
	return &Sender{
		api:     api,
		chatIDs: chatIDs,
	}
}

func (s *Sender) Send(text string) error {
	for _, id := range s.chatIDs {
		log.Println("TG SEND → chat:", id)
		msg := tgbotapi.NewMessage(id, text)
		if _, err := s.api.Send(msg); err != nil {
			log.Println("TG ERROR:", err)
		}
	}
	return nil
}
