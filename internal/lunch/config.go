package lunch

import (
	"time"

	"giautm.dev/com/internal/bot"
)

type Config struct {
	Port            string        `env:"PORT, default=8701"`
	ScheduleTimeout time.Duration `env:"SCHEDULE_TIMEOUT, default=5m"`

	Bot bot.Config
}

func (c *Config) BotConfig() *bot.Config {
	return &c.Bot
}
