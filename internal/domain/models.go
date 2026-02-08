package domain

type AvitoWebhook struct {
	Type    string         `json:"type"`
	Payload map[string]any `json:"payload"`
}
