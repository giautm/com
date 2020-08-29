package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"giautm.dev/com/pkg/logging"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	tgbotapi *tgbotapi.BotAPI
}

func NewTelegramBot(ctx context.Context, cfg *Config) (Bot, error) {
	logger := logging.FromContext(ctx)

	tele, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		return nil, err
	}

	logger.Infof("Authorized on account: %s", tele.Self.UserName)
	tele.Debug = cfg.DebugMode

	_, err = tele.Request(tgbotapi.NewWebhook(cfg.WebhookURL()))
	if err != nil {
		return nil, err
	}

	return &TelegramBot{
		tgbotapi: tele,
	}, nil
}

func (b TelegramBot) WebhookPath() string {
	return "/" + b.tgbotapi.Token
}

func (b *TelegramBot) UpdatesHandler(ctx context.Context, f func() error) http.HandlerFunc {
	logger := logging.FromContext(ctx)

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()

		var update tgbotapi.Update

		err := json.NewDecoder(req.Body).Decode(&update)
		if err == nil {
			err = f()
		}

		if err != nil {
			logger.Errorf("failed to process update: %v", err)

			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"status":"ERROR","message": "%v"}`, err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"status": "OK"}`)
	})
}

func (b *TelegramBot) ClosePoll(ctx context.Context, pollID string) error {
	return nil
}
