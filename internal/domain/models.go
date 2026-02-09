package domain

type AvitoEvent struct {
	Payload struct {
		Type  string `json:"type"`
		Value struct {
			UserID   int64 `json:"user_id"`
			AuthorID int64 `json:"author_id"`
			Content  struct {
				Text string `json:"text"`
			} `json:"content"`
		} `json:"value"`
	} `json:"payload"`
}
