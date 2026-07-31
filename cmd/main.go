package main

import (
	"context"
	"errors"
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
	startAvitoPoller(svc, accountAlias)

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

func startAvitoPoller(svc *domain.Service, accountAlias string) {
	configPath := os.Getenv("AVITO_SOURCES_CONFIG")
	accessToken := os.Getenv("AVITO_ACCESS_TOKEN")
	clientID := os.Getenv("AVITO_CLIENT_ID")
	clientSecret := os.Getenv("AVITO_CLIENT_SECRET")
	accountID := os.Getenv("AVITO_ACCOUNT_ID")
	chatID := os.Getenv("AVITO_CHAT_ID")

	sources, err := loadPollSources(configPath, accountAlias, accessToken, clientID, clientSecret, accountID, chatID)
	if err != nil {
		log.Fatal("AVITO SOURCES CONFIG ERROR:", err)
	}
	if len(sources) == 0 {
		log.Println("AVITO POLLER DISABLED: no enabled Avito sources")
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
	log.Println("AVITO_SOURCES_CONFIG:", configPath)
	log.Println("AVITO_SOURCES_COUNT:", len(sources))
	log.Println("POLL_INTERVAL:", interval)

	for _, source := range sources {
		log.Println("AVITO SOURCE ENABLED:", source.Name, source.Source, source.AccountID, source.ChatID)
		client := avito.NewClient(source.AccessToken, source.ClientID, source.ClientSecret)
		poller := avito.NewPoller(client, svc, source, interval)
		go poller.Run(context.Background())
	}
}

func loadPollSources(configPath, accountAlias, accessToken, clientID, clientSecret, accountID, chatID string) ([]avito.PollSource, error) {
	if configPath != "" {
		sources, err := avito.LoadSourcesConfig(configPath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return nil, err
			}
			return nil, err
		}
		return sources, nil
	}

	if accessToken == "" && (clientID == "" || clientSecret == "") {
		return nil, nil
	}
	if accountID == "" || chatID == "" {
		return nil, nil
	}

	return []avito.PollSource{
		{
			Name:         accountAlias,
			Source:       "polling: Проверка транспорта",
			AccountID:    accountID,
			ChatID:       chatID,
			AccessToken:  accessToken,
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Enabled:      true,
		},
	}, nil
}
