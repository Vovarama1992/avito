package delivery

import chi "github.com/go-chi/chi/v5"

func RegisterRoutes(r chi.Router, h *WebhookHandler) {
	r.Post("/webhook/avito", h.HandleAvitoWebhook)
}
