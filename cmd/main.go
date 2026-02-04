package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/Vovarama1992/avito/internal/delivery"
	"github.com/Vovarama1992/avito/internal/domain"
	"github.com/Vovarama1992/avito/internal/telegram"
)

func main() {
	log.Println("=== AVITO APP START ===")

	// --- env ---
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatIDStr := os.Getenv("TELEGRAM_CHAT_ID")
	port := os.Getenv("PORT")

	log.Println("PORT:", port)
	log.Println("TELEGRAM_CHAT_ID:", chatIDStr)
	if len(token) > 10 {
		log.Println("TELEGRAM_TOKEN_PREFIX:", token[:10])
	}

	// --- telegram ---
	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		log.Fatal("invalid TELEGRAM_CHAT_ID")
	}

	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Telegram bot authorized as:", api.Self.UserName)

	tgSender := telegram.NewSender(api, chatID)

	// --- domain ---
	svc := domain.NewService(tgSender)

	// --- http ---
	h := delivery.NewWebhookHandler(svc)
	r := chi.NewRouter()
	delivery.RegisterRoutes(r, h)

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
