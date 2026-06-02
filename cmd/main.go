package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"

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

	matrixSender := matrix.NewSender(
		matrixHomeserver,
		matrixAccessToken,
		matrixRoomID,
	)

	svc := domain.NewService(matrixSender)

	h := delivery.NewWebhookHandler(svc)
	r := chi.NewRouter()
	delivery.RegisterRoutes(r, h)

	addr := ":" + port
	log.Println("Listening on", addr)

	log.Fatal(http.ListenAndServe(addr, r))
}
