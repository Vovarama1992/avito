package delivery

import (
	"encoding/json"
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

	var payload domain.AvitoWebhookPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	h.svc.ProcessWebhook(r.Context(), payload)

	w.WriteHeader(http.StatusOK)
}
