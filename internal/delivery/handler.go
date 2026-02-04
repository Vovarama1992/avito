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

	var payload domain.AvitoWebhook

	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		log.Println("DECODE ERROR:", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if payload.Type != "message_new" {
		log.Println("SKIP EVENT TYPE:", payload.Type)
		w.WriteHeader(http.StatusOK)
		return
	}

	text := payload.Payload.Message.Text
	log.Println("TEXT:", text)

	h.svc.ProcessMessage(r.Context(), text)

	w.WriteHeader(http.StatusOK)
}
