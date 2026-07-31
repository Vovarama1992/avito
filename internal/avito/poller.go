package avito

import (
	"context"
	"errors"
	"time"

	"github.com/Vovarama1992/avito/internal/audit"
)

type SystemMessageHandler interface {
	ProcessSystemMessage(ctx context.Context, alias, source, accountID, chatID, text string)
}

type Poller struct {
	client            *Client
	handler           SystemMessageHandler
	name              string
	source            string
	accountID         string
	chatID            string
	interval          time.Duration
	startedAt         int64
	seen              map[string]struct{}
	ready             bool
	authAlert         bool
	consecutiveErrors int
	outageAlert       bool
}

func NewPoller(client *Client, handler SystemMessageHandler, source PollSource, interval time.Duration) *Poller {
	return &Poller{
		client:    client,
		handler:   handler,
		name:      source.Name,
		source:    source.Source,
		accountID: source.AccountID,
		chatID:    source.ChatID,
		interval:  interval,
		startedAt: time.Now().Unix(),
		seen:      make(map[string]struct{}),
	}
}

func (p *Poller) Run(ctx context.Context) {
	audit.Logf("AVITO POLLER STARTED: name=%s source=%s account_id=%s chat_id=%s", p.name, p.source, p.accountID, p.chatID)
	p.poll(ctx)

	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			audit.Logf("AVITO POLLER STOPPED: name=%s source=%s account_id=%s chat_id=%s", p.name, p.source, p.accountID, p.chatID)
			return
		case <-ticker.C:
			p.poll(ctx)
		}
	}
}

func (p *Poller) poll(ctx context.Context) {
	messages, err := p.client.GetMessages(ctx, p.accountID, p.chatID)
	if err != nil {
		p.consecutiveErrors++
		audit.Logf("AVITO POLL ERROR: account_id=%s chat_id=%s err=%v", p.accountID, p.chatID, err)
		if errors.Is(err, ErrUnauthorized) {
			p.sendAuthAlert(ctx)
		}
		p.sendOutageAlert(ctx, err)
		return
	}

	if p.outageAlert {
		p.sendRecoveryAlert(ctx)
	}
	p.consecutiveErrors = 0
	p.outageAlert = false
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
		if msg.Type != "system" {
			if p.ready {
				audit.Logf("AVITO POLLER RECEIVED BUT NOT SENT: reason=not_system id=%s type=%s direction=%s account_id=%s chat_id=%s", msg.ID, msg.Type, msg.Direction, p.accountID, p.chatID)
			}
			continue
		}
		if !p.ready && msg.Created < p.startedAt {
			continue
		}

		newCount++
		audit.Logf("AVITO POLLER NEW SYSTEM MESSAGE: id=%s direction=%s name=%s source=%s account_id=%s chat_id=%s text=%q", msg.ID, msg.Direction, p.name, p.source, p.accountID, p.chatID, preview(msg.Content.Text))
		p.handler.ProcessSystemMessage(ctx, p.name, p.source, p.accountID, p.chatID, msg.Content.Text)
	}

	if !p.ready {
		p.ready = true
		audit.Logf("AVITO POLLER BASELINE LOADED: account_id=%s chat_id=%s messages=%d system=%d new_system_sent=%d", p.accountID, p.chatID, len(messages), systemCount, newCount)
		return
	}

	audit.Logf("AVITO POLLER OK: account_id=%s chat_id=%s messages=%d system=%d new_system=%d", p.accountID, p.chatID, len(messages), systemCount, newCount)
}

func (p *Poller) sendAuthAlert(ctx context.Context) {
	if p.authAlert {
		return
	}
	p.authAlert = true
	audit.Logf("AVITO TOKEN ALERT SENT: account_id=%s chat_id=%s", p.accountID, p.chatID)

	p.handler.ProcessSystemMessage(
		ctx,
		p.name,
		"polling: ошибка Avito token",
		p.accountID,
		p.chatID,
		"Avito access token не работает или протух. Нужно обновить AVITO_ACCESS_TOKEN, иначе сообщения из \"Проверки транспорта\" не будут приходить.",
	)
}

func (p *Poller) sendOutageAlert(ctx context.Context, err error) {
	const threshold = 5
	if p.outageAlert || p.consecutiveErrors < threshold {
		return
	}
	p.outageAlert = true

	audit.Logf("AVITO API OUTAGE ALERT SENT: account_id=%s chat_id=%s consecutive_errors=%d", p.accountID, p.chatID, p.consecutiveErrors)
	p.handler.ProcessSystemMessage(
		ctx,
		p.name,
		"polling: ошибка Avito API",
		p.accountID,
		p.chatID,
		"Авария polling: Avito API не отвечает 5 запросов подряд. Сообщения из \"Проверки транспорта\" временно не читаются. Последняя ошибка: "+err.Error(),
	)
}

func (p *Poller) sendRecoveryAlert(ctx context.Context) {
	audit.Logf("AVITO API RECOVERY ALERT SENT: account_id=%s chat_id=%s", p.accountID, p.chatID)
	p.handler.ProcessSystemMessage(
		ctx,
		p.name,
		"polling: Avito API восстановился",
		p.accountID,
		p.chatID,
		"Polling снова работает: Avito API ответил, сообщения из \"Проверки транспорта\" снова читаются.",
	)
}

func preview(text string) string {
	const limit = 220
	if len(text) <= limit {
		return text
	}
	return text[:limit] + "..."
}
