package domain

type AvitoWebhook struct {
	ID        string `json:"id"`
	Version   string `json:"version"`
	Timestamp int64  `json:"timestamp"`
	Payload   struct {
		Type  string `json:"type"` // "message"
		Value struct {
			ID       string `json:"id"`
			ChatID   string `json:"chat_id"`
			UserID   int64  `json:"user_id"`
			AuthorID int64  `json:"author_id"`
			Created  int64  `json:"created"`
			Type     string `json:"type"` // "text", "image", ...
			ChatType string `json:"chat_type"`
			Content  struct {
				Text string `json:"text"`
			} `json:"content"`
			ItemID      int64  `json:"item_id"`
			PublishedAt string `json:"published_at"`
			FlowID      string `json:"flow_id"` // может быть пустым
		} `json:"value"`
	} `json:"payload"`
}
