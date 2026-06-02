package matrix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

type Sender struct {
	homeserver  string
	accessToken string
	roomID      string
	client      *http.Client
}

func NewSender(homeserver, accessToken, roomID string) *Sender {
	return &Sender{
		homeserver:  homeserver,
		accessToken: accessToken,
		roomID:      roomID,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *Sender) Send(text string) error {
	txnID := fmt.Sprintf("%d", time.Now().UnixNano())

	endpoint := fmt.Sprintf(
		"%s/_matrix/client/v3/rooms/%s/send/m.room.message/%s",
		s.homeserver,
		url.PathEscape(s.roomID),
		txnID,
	)

	body := map[string]string{
		"msgtype": "m.text",
		"body":    text,
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("matrix marshal message: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, endpoint, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("matrix create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.accessToken)
	req.Header.Set("Content-Type", "application/json")

	log.Println("MATRIX SEND → room:", s.roomID)

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("matrix send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("matrix bad status: %s", resp.Status)
	}

	return nil
}
