package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/Vovarama1992/avito/internal/delivery"
	"github.com/Vovarama1992/avito/internal/domain"
	"github.com/Vovarama1992/avito/internal/telegram"
)

func main() {
	log.Println("=== AVITO APP START ===")

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	port := os.Getenv("PORT")

	chatIDs := []int64{20461089, 6789440333}

	log.Println("PORT:", port)
	log.Println("CHAT IDS:", chatIDs)
	if len(token) > 10 {
		log.Println("TELEGRAM_TOKEN_PREFIX:", token[:10])
	}

	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Telegram bot authorized as:", api.Self.UserName)

	tgSender := telegram.NewSender(api, chatIDs)

	svc := domain.NewService(tgSender)

	h := delivery.NewWebhookHandler(svc)
	r := chi.NewRouter()
	delivery.RegisterRoutes(r, h)

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
