package domain

type AvitoWebhook struct {
	Type    string `json:"type"`
	Payload struct {
		Message struct {
			ID       string `json:"id"`
			Text     string `json:"text"`
			AuthorID string `json:"author_id"`
			ChatID   string `json:"chat_id"`
			Created  int64  `json:"created"`
		} `json:"message"`
	} `json:"payload"`
}
