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

	body, _ := io.ReadAll(r.Body)
	log.Println("=== AVITO WEBHOOK RAW ===")
	log.Println(string(body))

	var ev domain.AvitoEvent
	if err := json.Unmarshal(body, &ev); err != nil {
		log.Println("JSON ERROR:", err)
		w.WriteHeader(http.StatusOK)
		return
	}

	h.svc.ProcessEvent(r.Context(), ev)

	w.WriteHeader(http.StatusOK)
}
