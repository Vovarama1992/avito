package delivery

import (
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

	h.svc.ProcessMessage(r.Context(), string(body))

	w.WriteHeader(http.StatusOK)
}

func extractText(m map[string]any) string {
	// payload.message.content.text
	if p, ok := m["payload"].(map[string]any); ok {
		if msg, ok := p["message"].(map[string]any); ok {
			if content, ok := msg["content"].(map[string]any); ok {
				if t, ok := content["text"].(string); ok {
					return t
				}
			}
			// fallback: payload.message.text
			if t, ok := msg["text"].(string); ok {
				return t
			}
		}
	}
	return ""
}
