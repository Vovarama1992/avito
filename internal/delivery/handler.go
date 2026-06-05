package delivery

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/Vovarama1992/avito/internal/domain"
)

type WebhookHandler struct {
	svc *domain.Service
}

func NewWebhookHandler(svc *domain.Service) *WebhookHandler {
	return &WebhookHandler{svc: svc}
}

type AvitoEventWebhook struct {
	Event string `json:"event"`
	Data  struct {
		MessageID string `json:"message_id"`
		Text      string `json:"text"`
	} `json:"data"`
}

func (h *WebhookHandler) HandleAvitoWebhook(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, _ := io.ReadAll(r.Body)

	log.Println("=== AVITO WEBHOOK RAW ===")
	log.Println(string(body))

	payloadBody, err := extractPayloadBody(body)
	if err != nil {
		log.Println("PAYLOAD EXTRACT ERROR:", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if isPing(payloadBody) {
		log.Println("PING OK")
		w.WriteHeader(http.StatusOK)
		return
	}

	evt, ok := parseWebhook(payloadBody)
	if !ok {
		log.Println("SKIP: unsupported webhook format")
		w.WriteHeader(http.StatusOK)
		return
	}

	log.Println("PAYLOAD TYPE:", evt.Payload.Type)
	log.Println("MSG TYPE:", evt.Payload.Value.Type)

	h.svc.ProcessWebhook(r.Context(), evt)

	w.WriteHeader(http.StatusOK)
}

func extractPayloadBody(body []byte) ([]byte, error) {
	raw := strings.TrimSpace(string(body))

	if strings.HasPrefix(raw, "{") {
		return []byte(raw), nil
	}

	values, err := url.ParseQuery(raw)
	if err != nil {
		return nil, err
	}

	payload := values.Get("payload")
	if payload == "" {
		return []byte(raw), nil
	}

	return []byte(payload), nil
}

func isPing(body []byte) bool {
	var raw map[string]any
	if err := json.Unmarshal(body, &raw); err != nil {
		return false
	}

	if raw["event"] == "ping" {
		return true
	}

	if raw["test"] == "ping" {
		return true
	}

	return false
}

func parseWebhook(body []byte) (domain.AvitoWebhook, bool) {
	var evt domain.AvitoWebhook
	if err := json.Unmarshal(body, &evt); err == nil {
		if evt.Payload.Value.Content.Text != "" {
			return evt, true
		}
	}

	var avitoEvt AvitoEventWebhook
	if err := json.Unmarshal(body, &avitoEvt); err != nil {
		log.Println("DECODE ERROR:", err)
		return domain.AvitoWebhook{}, false
	}

	if avitoEvt.Event != "message_new" {
		log.Println("SKIP: unsupported event:", avitoEvt.Event)
		return domain.AvitoWebhook{}, false
	}

	if avitoEvt.Data.Text == "" {
		log.Println("SKIP: empty text")
		return domain.AvitoWebhook{}, false
	}

	evt.Payload.Type = "message"
	evt.Payload.Value.Type = "text"
	evt.Payload.Value.Content.Text = avitoEvt.Data.Text

	return evt, true
}
