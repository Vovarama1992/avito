package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Sender struct {
	api    *tgbotapi.BotAPI
	chatID int64
}

func NewSender(api *tgbotapi.BotAPI, chatID int64) *Sender {
	log.Println("TELEGRAM SENDER INIT. CHAT ID:", chatID)
	return &Sender{
		api:    api,
		chatID: chatID,
	}
}

func (s *Sender) Send(text string) error {
	log.Println("TG SEND → chat:", s.chatID)
	log.Println("TG TEXT:", text)

	msg := tgbotapi.NewMessage(s.chatID, text)
	res, err := s.api.Send(msg)
	if err != nil {
		log.Println("TG ERROR:", err)
		return err
	}

	log.Println("TG OK. MessageID:", res.MessageID)
	return nil
}
