package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"giautm.dev/com/pkg/logging"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const maxPollOptions = 10

type TelegramBot struct {
	bot  *tgbotapi.BotAPI
	repo BotRepo
}

func NewTelegramBot(ctx context.Context, cfg *Config) (BotWebhook, error) {
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
		bot:  tele,
		repo: NewMemoryBotRepo(),
	}, nil
}

func (b TelegramBot) WebhookPath() string {
	return "/" + b.bot.Token
}

func (b *TelegramBot) UpdatesHandler(ctx context.Context, h Handler) http.HandlerFunc {
	logger := logging.FromContext(ctx)

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()

		var update tgbotapi.Update

		err := json.NewDecoder(req.Body).Decode(&update)
		if err == nil {
			if update.PollAnswer != nil {
				p := update.PollAnswer
				err = h.OnPollAnswer(&TelegramChatContext{
					bot: b,
					ctx: req.Context(),
				}, p.PollID, p.OptionIDs, p.User.ID)
			}

			if update.Message != nil {
				msg := update.Message
				if cmd := msg.Command(); cmd != "" {
					err = h.OnCommand(&TelegramChatContext{
						bot: b,
						ctx: req.Context(),
						msg: msg,
					}, cmd, msg.CommandArguments())
				}
			}
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

func (b *TelegramBot) DeleteMessage(ctx context.Context, chatID int64, msgID int) error {
	reply := tgbotapi.DeleteMessageConfig{}
	reply.ChatID = chatID
	reply.MessageID = msgID

	_, err := b.bot.Send(reply)
	return err
}

func (b *TelegramBot) SendMessage(ctx context.Context, chatID int64, msg string) error {
	reply := tgbotapi.NewMessage(chatID, msg)
	// reply.ReplyToMessageID = t.msg.MessageID

	_, err := b.bot.Send(reply)
	return err
}

func (b *TelegramBot) SendPoll(ctx context.Context, chatID int64, question string, options []string) (string, error) {
	reply := tgbotapi.NewPoll(chatID, question, options...)
	reply.IsAnonymous = false
	// reply.ReplyToMessageID = t.msg.MessageID

	// TODO: Update tgbotapi to use context
	msgPoll, err := b.bot.Send(reply)
	if err != nil {
		return "", err
	}

	pollID := msgPoll.Poll.ID
	err = b.repo.SavePoll(ctx, pollID, msgPoll.MessageID)
	if err != nil {
		return "", err
	}

	return pollID, nil
}

func (b *TelegramBot) StopPoll(ctx context.Context, chatID int64, pollID string) error {
	msgID, err := b.repo.FetchMessageID(ctx, pollID)
	if err != nil {
		return err
	}

	reply := tgbotapi.NewStopPoll(chatID, msgID)
	_, err = b.bot.Send(reply)
	return err
}

func chunkBy(items []string, chunkSize int) (chunks [][]string) {
	for chunkSize < len(items) {
		items, chunks = items[chunkSize:], append(chunks, items[0:chunkSize:chunkSize])
	}

	return append(chunks, items)
}
