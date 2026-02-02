package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Sender struct {
	api    *tgbotapi.BotAPI
	chatID int64
}

func NewSender(api *tgbotapi.BotAPI, chatID int64) *Sender {
	return &Sender{
		api:    api,
		chatID: chatID,
	}
}

func (s *Sender) Send(text string) error {
	msg := tgbotapi.NewMessage(s.chatID, text)
	_, err := s.api.Send(msg)
	return err
}
