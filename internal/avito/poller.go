package avito

import (
	"context"
	"errors"
	"log"
	"time"
)

type SystemMessageHandler interface {
	ProcessSystemMessage(ctx context.Context, source, accountID, chatID, text string)
}

type Poller struct {
	client    *Client
	handler   SystemMessageHandler
	accountID string
	chatID    string
	interval  time.Duration
	seen      map[string]struct{}
	ready     bool
	authAlert bool
}

func NewPoller(client *Client, handler SystemMessageHandler, accountID, chatID string, interval time.Duration) *Poller {
	return &Poller{
		client:    client,
		handler:   handler,
		accountID: accountID,
		chatID:    chatID,
		interval:  interval,
		seen:      make(map[string]struct{}),
	}
}

func (p *Poller) Run(ctx context.Context) {
	log.Println("AVITO POLLER STARTED")
	p.poll(ctx)

	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("AVITO POLLER STOPPED")
			return
		case <-ticker.C:
			p.poll(ctx)
		}
	}
}

func (p *Poller) poll(ctx context.Context) {
	messages, err := p.client.GetMessages(ctx, p.accountID, p.chatID)
	if err != nil {
		log.Println("AVITO POLL ERROR:", err)
		if errors.Is(err, ErrUnauthorized) {
			p.sendAuthAlert(ctx)
		}
		return
	}

	p.authAlert = false
	systemCount := 0
	newCount := 0

	for _, msg := range messages {
		if msg.ID == "" {
			continue
		}
		if msg.Type == "system" {
			systemCount++
		}
		if _, ok := p.seen[msg.ID]; ok {
			continue
		}

		p.seen[msg.ID] = struct{}{}
		if !p.ready {
			continue
		}
		if msg.Type != "system" {
			log.Printf("AVITO POLLER RECEIVED BUT NOT SENT: reason=not_system id=%s type=%s direction=%s", msg.ID, msg.Type, msg.Direction)
			continue
		}

		newCount++
		p.handler.ProcessSystemMessage(ctx, "polling: Проверка транспорта", p.accountID, p.chatID, msg.Content.Text)
	}

	if !p.ready {
		p.ready = true
		log.Printf("AVITO POLLER BASELINE LOADED: messages=%d system=%d", len(messages), systemCount)
		return
	}

	log.Printf("AVITO POLLER OK: messages=%d system=%d new_system=%d", len(messages), systemCount, newCount)
}

func (p *Poller) sendAuthAlert(ctx context.Context) {
	if p.authAlert {
		return
	}
	p.authAlert = true

	p.handler.ProcessSystemMessage(
		ctx,
		"polling: ошибка Avito token",
		p.accountID,
		p.chatID,
		"Avito access token не работает или протух. Нужно обновить AVITO_ACCESS_TOKEN, иначе сообщения из \"Проверки транспорта\" не будут приходить.",
	)
}
