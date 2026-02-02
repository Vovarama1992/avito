package domain

type AvitoWebhookPayload struct {
	Messages []AvitoMessage `json:"messages"`
}

type AvitoMessage struct {
	Type    string       `json:"type"`
	Content AvitoContent `json:"content"`
}

type AvitoContent struct {
	Text string `json:"text"`
}
