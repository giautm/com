package lunch

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"giautm.dev/com/internal/bot"
	"giautm.dev/com/internal/serverenv"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type PollSender interface {
	Chunks(options []string) [][]string
	SendPoll(ctx context.Context, question string, options []string) (messageId int, pollId string, err error)
	SendMessage(ctx context.Context, message string) error
}

type Service interface {
	NewLunch(ctx context.Context,
		chatID int64,
		msg string,
		timestamp time.Time,
		sender PollSender,
	) error
}

type Poll struct {
}

type PollRepo interface {
	Save(ctx context.Context, poll *Poll)
}

type Subscribe struct {
}

func (s *Subscribe) SubscribeForChat(chatID int, subscribe bool) {

}

type SubscribeRepo interface {
	Update(ctx context.Context, userID int, updater func(sub *Subscribe) error) error
}
type LunchHandler struct {
	pollRepo PollRepo
	subRepo  SubscribeRepo
	cfg      *Config
	bot      bot.Bot
}

func NewHandler(cfg *Config, env *serverenv.ServerEnv) (*LunchHandler, error) {
	return &LunchHandler{
		bot: env.Bot(),
		cfg: cfg,
	}, nil
}

func (s *LunchHandler) parseOptions(text string) (options []string, leadtime time.Duration) {
	leadtime = 3 * time.Hour

	firstNonEmpty := true
	lines := strings.Split(text, "\n")
	unique := make(map[string]struct{})
	for _, line := range lines {
		option := strings.TrimSpace(strings.TrimRight(line, ".…"))
		if option != "" {
			if firstNonEmpty {
				firstNonEmpty = false
				if tmp, err := strconv.ParseInt(option, 10, 32); err == nil {
					leadtime = time.Duration(tmp) * time.Hour
					log.Printf("duration: %d\n", tmp)
					continue
				}
			}

			if _, ok := unique[option]; !ok {
				options = append(options, option)

				// Ensure option is unique
				unique[option] = struct{}{}
			}
		}
	}

	return options, leadtime
}

func TimeIn(t time.Time) (time.Time, error) {
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err == nil {
		t = t.In(loc)
	}
	return t, err
}

func (s *LunchHandler) NewLunch(ctx context.Context,
	chatID int64,
	msg string,
	timestamp time.Time,
	sender PollSender,
) error {

	options, hours := s.parseOptions(msg)
	if hours <= 0 {
		sender.SendMessage(ctx, "Thời gian không hợp lệ, cần lớn hơn 0")
		return nil
	}

	if len(options) > 0 {
		polls := sender.Chunks(options)
		totalPolls := len(polls)

		// Because Telegram Bot only support 10 minutes,
		// so We need close poll manually
		//
		// https://core.telegram.org/bots/api#sendpoll
		closeDate, _ := TimeIn(timestamp.Add(hours))
		closeDateStr := closeDate.Format("15:04 02/01")

		for idx, pollOptions := range polls {
			question := fmt.Sprintf("%s Trưa nay ăn gì?", timestamp.Format("2006-01-02"))
			if totalPolls > 1 {
				question = fmt.Sprintf("[%d/%d] %s", idx+1, totalPolls, question)
			}
			question = fmt.Sprintf("%s\n\nChốt cơm lúc: %s", question, closeDateStr)

			msgID, pollID, err := sender.SendPoll(ctx, question, pollOptions)
			if err != nil {
				return err
			}

			log.Printf("MessageID: %d, PollID: %s", msgID, pollID)
		}
	}

	return nil
}

type ChatContext interface {
	Context() context.Context

	DeleteIncomeMessage() error
	DeleteMessage(msgID int) error
	ReplyMessage(message string) (msgID int, err error)
	SendMessage(message string) (msgID int, err error)

	ClosePoll(pollID string) error
	SendPoll(question string, options []string) (pollID string, msgID int, err error)
}

func (s *LunchHandler) process(ctx context.Context, update *tgbotapi.Update) error {
	if update.PollAnswer != nil {
		// poll := update.PollAnswer

		// s.UpdatePollAnswer(ctx, poll.PollID, func(p *Poll) error {

		// })
	}

	if update.Message != nil {
		// chat := update.Message.Chat

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := update.Message
		cmd := msg.Command()
		switch cmd {
		case "start":
			if msg.Chat.IsPrivate() {
				return nil
			}
			// admins, err := s.bot.GetChatAdministrators(tgbotapi.ChatAdministratorsConfig{
			// 	ChatConfig: tgbotapi.ChatConfig{
			// 		ChatID: msg.Chat.ID,
			// 	},
			// })
			// if err != nil {
			// 	log.Printf("failed to get admins: %v", err)
			// 	return err
			// }

			// for _, admin := range admins {
			// 	if admin.CustomTitle == "chu-no" {
			// 	}
			// 	if admin.CustomTitle == "order" {
			// 	}
			// }
			break
		// case "lunch":
		// 	arg := msg.CommandArguments()
		// 	timestamp := time.Unix(int64(msg.Date), 0)

		// 	return s.NewLunch(ctx, chat.ID, arg, timestamp, tgbot.NewPoll(s.bot, msg))
		case "subscribe":

			return nil
		case "unsubscribe":
			return nil
		}
	}

	return nil
}

func (s *LunchHandler) Subscribe(chat ChatContext) error {
	err := s.subRepo.Update(chat.Context(), 12, func(sub *Subscribe) error {
		sub.SubscribeForChat(12, true)
		return nil
	})
	if err == nil {
		_, err = chat.ReplyMessage(`Đã thêm bạn vào danh sách thông báo thành công!
Để huỷ nhắc cơm, gửi lệnh /unsubscribe`)
	}

	return err
}
