package main

import (
	"net/http"
	"os"

	"giautm.dev/com/internal/lunch"
	"giautm.dev/com/internal/tgbot"
)

func main() {
	port := os.Getenv("PORT")
	botCfg := &tgbot.Config{
		BaseURL:  os.Getenv("CLOUD_RUN_SERVICE_URL"),
		BotToken: os.Getenv("BOT_TOKEN"),
	}
	bot, err := tgbot.NewBot(botCfg)
	if err != nil {
		panic(err)
	}

	s := lunch.NewHandler(&lunch.Config{}, bot)

	mux := http.NewServeMux()
	mux.Handle(botCfg.WebhookPath(), s.Handle())

	err = http.ListenAndServe(":"+port, mux)
	if err != nil {
		panic(err)
	}
}
