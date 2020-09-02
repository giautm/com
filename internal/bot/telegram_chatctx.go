package bot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramChatContext struct {
	ctx context.Context
	bot Bot
	msg *tgbotapi.Message
}

var _ ContextChat = &TelegramChatContext{}

func (c *TelegramChatContext) Bot() Bot {
	return c.bot
}

func (c *TelegramChatContext) Context() context.Context {
	return c.ctx
}

func (c *TelegramChatContext) DeleteIncomeMessage(ctx context.Context) error {
	return c.bot.DeleteMessage(ctx, c.msg.Chat.ID, c.msg.MessageID)
}

func (c *TelegramChatContext) SendMessage(ctx context.Context, message string) error {
	return c.bot.SendMessage(ctx, c.msg.Chat.ID, message)
}

func (c *TelegramChatContext) SendPoll(ctx context.Context, question string, options []string) (pollID string, err error) {
	return c.bot.SendPoll(ctx, c.msg.Chat.ID, question, options)
}
