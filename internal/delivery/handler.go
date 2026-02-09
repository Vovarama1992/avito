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

	var evt domain.AvitoWebhook
	if err := json.Unmarshal(body, &evt); err != nil {
		log.Println("DECODE ERROR:", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	log.Println("PAYLOAD TYPE:", evt.Payload.Type) // обычно "message"
	log.Println("MSG TYPE:", evt.Payload.Value.Type)

	h.svc.ProcessWebhook(r.Context(), evt)

	w.WriteHeader(http.StatusOK)
}
