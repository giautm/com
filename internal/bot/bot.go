package bot

import (
	"context"
	"fmt"
	"net/http"
)

type Bot interface {
	UpdatesHandler(ctx context.Context, update func() error) http.HandlerFunc
	WebhookPath() string

	ClosePoll(ctx context.Context, pollID string) error
}

// SetupBotFor returns the bot for the given type, or an error
// if one does not exist.
func SetupBotFor(ctx context.Context, cfg *Config) (Bot, error) {
	typ := cfg.BotType
	switch typ {
	case BotTypeTelegram:
		return NewTelegramBot(ctx, cfg)
	}

	return nil, fmt.Errorf("unknown bot type: %v", typ)
}
