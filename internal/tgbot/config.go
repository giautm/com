package tgbot

type Config struct {
	BaseURL  string
	BotToken string
}

func (c Config) WebhookPath() string {
	return "/" + c.BotToken
}

func (c Config) WebhookURL() string {
	return c.BaseURL + c.WebhookPath()
}
