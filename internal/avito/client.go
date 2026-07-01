package avito

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

var ErrUnauthorized = errors.New("avito unauthorized")

type Client struct {
	baseURL     string
	accessToken string
	httpClient  *http.Client
}

type Message struct {
	ID        string `json:"id"`
	AuthorID  int64  `json:"author_id"`
	Created   int64  `json:"created"`
	Type      string `json:"type"`
	Direction string `json:"direction"`
	Content   struct {
		Text string `json:"text"`
	} `json:"content"`
}

type messagesResponse struct {
	Messages []Message `json:"messages"`
	Meta     struct {
		HasMore bool `json:"has_more"`
	} `json:"meta"`
}

func NewClient(accessToken string) *Client {
	return &Client{
		baseURL:     "https://api.avito.ru",
		accessToken: accessToken,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (c *Client) GetMessages(ctx context.Context, accountID, chatID string) ([]Message, error) {
	endpoint := fmt.Sprintf(
		"%s/messenger/v3/accounts/%s/chats/%s/messages/",
		c.baseURL,
		url.PathEscape(accountID),
		url.PathEscape(chatID),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("avito create messages request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("avito get messages: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, fmt.Errorf("%w: %s", ErrUnauthorized, resp.Status)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("avito get messages bad status: %s", resp.Status)
	}

	var out messagesResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("avito decode messages: %w", err)
	}

	return out.Messages, nil
}
