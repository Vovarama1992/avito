package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/Vovarama1992/avito/internal/avito"
	"github.com/Vovarama1992/avito/internal/delivery"
	"github.com/Vovarama1992/avito/internal/domain"
	"github.com/Vovarama1992/avito/internal/matrix"
)

func main() {
	log.Println("=== AVITO APP START ===")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	accountAlias := os.Getenv("ACCOUNT_ALIAS")
	if accountAlias == "" {
		accountAlias = "Самара Jaecoo"
	}

	matrixHomeserver := os.Getenv("MATRIX_HOMESERVER")
	matrixAccessToken := os.Getenv("MATRIX_ACCESS_TOKEN")
	matrixRoomID := os.Getenv("MATRIX_ROOM_ID")

	if matrixHomeserver == "" {
		log.Fatal("MATRIX_HOMESERVER is required")
	}
	if matrixAccessToken == "" {
		log.Fatal("MATRIX_ACCESS_TOKEN is required")
	}
	if matrixRoomID == "" {
		log.Fatal("MATRIX_ROOM_ID is required")
	}

	log.Println("PORT:", port)
	log.Println("MATRIX_HOMESERVER:", matrixHomeserver)
	log.Println("MATRIX_ROOM_ID:", matrixRoomID)
	log.Println("ACCOUNT_ALIAS:", accountAlias)
	log.Println("WEBHOOK_FORWARD_ENABLED:", webhookForwardEnabled())

	matrixSender := matrix.NewSender(
		matrixHomeserver,
		matrixAccessToken,
		matrixRoomID,
	)

	svc := domain.NewService(matrixSender, accountAlias, webhookForwardEnabled())
	startAvitoPoller(svc)

	h := delivery.NewWebhookHandler(svc)
	r := chi.NewRouter()
	delivery.RegisterRoutes(r, h)

	addr := ":" + port
	log.Println("Listening on", addr)

	log.Fatal(http.ListenAndServe(addr, r))
}

func webhookForwardEnabled() bool {
	raw := os.Getenv("WEBHOOK_FORWARD_ENABLED")
	return raw == "" || raw == "true" || raw == "1" || raw == "yes"
}

func startAvitoPoller(svc *domain.Service) {
	accessToken := os.Getenv("AVITO_ACCESS_TOKEN")
	accountID := os.Getenv("AVITO_ACCOUNT_ID")
	chatID := os.Getenv("AVITO_CHAT_ID")

	if accessToken == "" || accountID == "" || chatID == "" {
		log.Println("AVITO POLLER DISABLED: AVITO_ACCESS_TOKEN, AVITO_ACCOUNT_ID or AVITO_CHAT_ID is empty")
		return
	}

	interval := 30 * time.Second
	if raw := os.Getenv("POLL_INTERVAL_SECONDS"); raw != "" {
		seconds, err := strconv.Atoi(raw)
		if err != nil || seconds <= 0 {
			log.Println("INVALID POLL_INTERVAL_SECONDS, using default 30")
		} else {
			interval = time.Duration(seconds) * time.Second
		}
	}

	log.Println("AVITO POLLER ENABLED")
	log.Println("AVITO_ACCOUNT_ID:", accountID)
	log.Println("AVITO_CHAT_ID:", chatID)
	log.Println("POLL_INTERVAL:", interval)

	client := avito.NewClient(accessToken)
	poller := avito.NewPoller(client, svc, accountID, chatID, interval)
	go poller.Run(context.Background())
}
