package avito

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/Vovarama1992/avito/internal/audit"
)

var ErrUnauthorized = errors.New("avito unauthorized")

type Client struct {
	baseURL      string
	clientID     string
	clientSecret string
	accessToken  string
	httpClient   *http.Client
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

type tokenResponse struct {
	AccessToken string `json:"access_token"`
}

func NewClient(accessToken, clientID, clientSecret string) *Client {
	return &Client{
		baseURL:      "https://api.avito.ru",
		clientID:     clientID,
		clientSecret: clientSecret,
		accessToken:  accessToken,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (c *Client) GetMessages(ctx context.Context, accountID, chatID string) ([]Message, error) {
	if c.accessToken == "" {
		if err := c.RefreshToken(ctx); err != nil {
			return nil, err
		}
	}

	messages, err := c.getMessages(ctx, accountID, chatID)
	if err != nil && c.canRefreshToken() {
		audit.Logf("AVITO GET MESSAGES ERROR: refreshing token before retry err=%v", err)
		if refreshErr := c.RefreshToken(ctx); refreshErr != nil {
			if errors.Is(err, ErrUnauthorized) {
				return nil, refreshErr
			}
			return nil, err
		}
		return c.getMessages(ctx, accountID, chatID)
	}

	return messages, err
}

func (c *Client) RefreshToken(ctx context.Context) error {
	if !c.canRefreshToken() {
		return ErrUnauthorized
	}

	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	form.Set("client_id", c.clientID)
	form.Set("client_secret", c.clientSecret)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/token/", bytes.NewBufferString(form.Encode()))
	if err != nil {
		return fmt.Errorf("avito create token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("avito get token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return fmt.Errorf("%w: token endpoint %s", ErrUnauthorized, resp.Status)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("avito get token bad status: %s", resp.Status)
	}

	var out tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return fmt.Errorf("avito decode token: %w", err)
	}
	if out.AccessToken == "" {
		return errors.New("avito token response is empty")
	}

	c.accessToken = out.AccessToken
	audit.Logf("AVITO TOKEN REFRESH OK")
	return nil
}

func (c *Client) canRefreshToken() bool {
	return c.clientID != "" && c.clientSecret != ""
}

func (c *Client) getMessages(ctx context.Context, accountID, chatID string) ([]Message, error) {
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
