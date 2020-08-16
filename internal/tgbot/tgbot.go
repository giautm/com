package tgbot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func NewBot(cfg *Config) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		return nil, err
	}
	log.Printf("Authorized on account: %s", bot.Self.UserName)
	bot.Debug = true

	_, err = bot.Request(tgbotapi.NewWebhook(cfg.WebhookURL()))
	if err != nil {
		return nil, err
	}

	info, err := bot.GetWebhookInfo()
	if err != nil {
		return nil, err
	}
	if info.LastErrorDate != 0 {
		log.Printf("failed to set webhook: %s", info.LastErrorMessage)
	}

	return bot, nil
}
