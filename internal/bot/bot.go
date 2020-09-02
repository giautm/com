package bot

import (
	"context"
	"fmt"
	"net/http"
)

type Context interface {
	Bot() Bot
	Context() context.Context
}

type ContextChat interface {
	Context

	DeleteIncomeMessage(context.Context) error
	SendMessage(ctx context.Context, message string) error
	SendPoll(ctx context.Context, question string, options []string) (pollID string, err error)
}

type Bot interface {
	DeleteMessage(ctx context.Context, chatID int64, msgID int) error
	SendMessage(ctx context.Context, chatID int64, message string) error
	SendPoll(ctx context.Context, chatID int64, question string, options []string) (string, error)
	StopPoll(ctx context.Context, chatID int64, pollID string) error
}

type Handler interface {
	OnCommand(ctx ContextChat, cmd string, arguments string) error
	OnPollAnswer(ctx Context, pollID string, options []int, userID int) error
}

type Webhook interface {
	UpdatesHandler(ctx context.Context, h Handler) http.HandlerFunc
	WebhookPath() string
}

type BotWebhook interface {
	Bot
	Webhook
}

// SetupBotFor returns the bot for the given type, or an error
// if one does not exist.
func SetupBotFor(ctx context.Context, cfg *Config) (BotWebhook, error) {
	typ := cfg.BotType
	switch typ {
	case BotTypeTelegram:
		return NewTelegramBot(ctx, cfg)
	}

	return nil, fmt.Errorf("unknown bot type: %v", typ)
}
