package bot

import (
	"context"
	"fmt"
	"net/http"
)

type Bot interface {
	WebhookPath() string
	UpdatesHandler(ctx context.Context, update func() error) http.HandlerFunc

	ClosePoll(ctx context.Context, pollID string) error
}

// BotFor returns the secret manager for the given type, or an error
// if one does not exist.
func BotFor(ctx context.Context, cfg *Config) (Bot, error) {
	typ := cfg.BotType
	switch typ {
	case BotTypeTelegram:
		return NewTelegramBot(ctx, cfg)
	}

	return nil, fmt.Errorf("unknown bot type: %v", typ)
}

// func NewBot(cfg *Config) (*tgbotapi.BotAPI, error) {
// 	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
// 	if err != nil {
// 		return nil, err
// 	}
// 	log.Printf("Authorized on account: %s", bot.Self.UserName)
// 	bot.Debug = true

// 	_, err = bot.Request(tgbotapi.NewWebhook(cfg.WebhookURL()))
// 	if err != nil {
// 		return nil, err
// 	}

// 	info, err := bot.GetWebhookInfo()
// 	if err != nil {
// 		return nil, err
// 	}
// 	if info.LastErrorDate != 0 {
// 		log.Printf("failed to set webhook: %s", info.LastErrorMessage)
// 	}

// 	return bot, nil
// }
