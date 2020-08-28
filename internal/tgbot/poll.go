package tgbot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const maxPollOptions = 10

type telegramPoll struct {
	bot *tgbotapi.BotAPI
	msg *tgbotapi.Message
}

func NewPoll(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) *telegramPoll {
	return &telegramPoll{
		bot: bot,
		msg: msg,
	}
}

func (t telegramPoll) Chunks(options []string) [][]string {
	return chunkBy(options, maxPollOptions)
}

func (t *telegramPoll) SendMessage(ctx context.Context, msg string) error {
	reply := tgbotapi.NewMessage(t.msg.Chat.ID, msg)
	reply.ReplyToMessageID = t.msg.MessageID

	_, err := t.bot.Send(reply)
	return err
}

func (t *telegramPoll) SendPoll(ctx context.Context, question string, options []string) (int, string, error) {
	chatID := ChatFromContext(ctx)

	reply := tgbotapi.NewPoll(chatID, question, options...)
	reply.IsAnonymous = false
	reply.ReplyToMessageID = t.msg.MessageID

	// TODO: Update tgbotapi to use context
	msgPoll, err := t.bot.Send(reply)
	if err != nil {
		return -1, "", err
	}

	return msgPoll.MessageID, msgPoll.Poll.ID, nil
}
