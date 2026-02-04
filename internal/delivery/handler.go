package delivery

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/Vovarama1992/avito/internal/domain"
)

type WebhookHandler struct {
	svc *domain.Service
}

func NewWebhookHandler(svc *domain.Service) *WebhookHandler {
	return &WebhookHandler{svc: svc}
}

func (h *WebhookHandler) HandleAvitoWebhook(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	log.Println("=== WEBHOOK HIT ===")
	log.Println("METHOD:", r.Method, "PATH:", r.URL.Path)
	log.Println("REMOTE:", r.RemoteAddr)

	bodyBytes, _ := io.ReadAll(r.Body)
	log.Println("RAW BODY:", string(bodyBytes))

	var payload domain.AvitoWebhookPayload
	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		log.Println("DECODE ERROR:", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	log.Printf("MESSAGES COUNT: %d\n", len(payload.Messages))

	h.svc.ProcessWebhook(r.Context(), payload)

	w.WriteHeader(http.StatusOK)
}
