package bot

// BotType represents a type of secret manager.
type BotType string

const (
	BotTypeTelegram BotType = "TELEGRAM"
)

type Config struct {
	BotType   BotType `env:"BOT_TYPE, default=TELEGRAM"`
	BaseURL   string  `env:"WEBHOOK_BASE_URL"`
	BotToken  string  `env:"BOT_TOKEN"`
	DebugMode bool    `env:"BOT_DEBUG, default=false"`
}

func (c Config) WebhookPath() string {
	return "/" + c.BotToken
}

func (c Config) WebhookURL() string {
	return c.BaseURL + c.WebhookPath()
}
