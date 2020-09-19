module giautm.dev/com

go 1.13

require (
	github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.0.0-rc1
	github.com/sethvargo/go-envconfig v0.3.1
	github.com/sethvargo/go-signalcontext v0.1.0
	github.com/stretchr/testify v1.6.1 // indirect
	go.opencensus.io v0.22.4 // indirect
	go.uber.org/zap v1.15.0
	gocloud.dev v0.20.0
	golang.org/x/tools v0.0.0-20200615222825-6aa8f57aacd9 // indirect
)

replace github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.0.0-rc1 => github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.0.0-rc1.0.20200723221353-2f7211a7085f
