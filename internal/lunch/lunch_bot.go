package lunch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"giautm.dev/com/internal/lunch/domain"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type PollSender interface {
	Chunks(options []string) [][]string
	SendPoll(ctx context.Context, question string, options []string) (messageId int, pollId string, err error)
	SendMessage(ctx context.Context, message string) error
}

type LunchHandler struct {
	groupRepo domain.GroupRepo
	cfg       *Config
	bot       *tgbotapi.BotAPI
}

func NewHandler(cfg *Config, bot *tgbotapi.BotAPI) *LunchHandler {
	return &LunchHandler{
		bot:       bot,
		groupRepo: NewMemoryGroupRepo(),
		cfg:       cfg,
	}
}

func TimeIn(t time.Time) (time.Time, error) {
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err == nil {
		t = t.In(loc)
	}
	return t, err
}

func (s *LunchHandler) process(ctx context.Context, update *tgbotapi.Update) error {
	if update.PollAnswer != nil {
		// poll := update.PollAnswer

		// s.UpdatePollAnswer(ctx, poll.PollID, func(p *Poll) error {

		// })
	}

	if update.Message != nil {
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := update.Message
		cmd := msg.Command()
		switch cmd {
		case "start":
			return s.handleStart(ctx, msg)
		case "lunch":
			return s.handleLunch(ctx, msg)
		}
	}

	return nil
}

func (s *LunchHandler) handleStart(ctx context.Context, msg *tgbotapi.Message) error {
	if msg.Chat.IsPrivate() {
		return nil
	}
	admins, err := s.bot.GetChatAdministrators(tgbotapi.ChatAdministratorsConfig{
		ChatConfig: tgbotapi.ChatConfig{
			ChatID: msg.Chat.ID,
		},
	})
	if err != nil {
		log.Printf("failed to get admins: %v", err)
		return err
	}

	for _, admin := range admins {
		if admin.CustomTitle == "chu-no" {
		}
		if admin.CustomTitle == "order" {
		}
	}
	return nil
}

func (s *LunchHandler) handleLunch(ctx context.Context, msg *tgbotapi.Message) error {
	arg := msg.CommandArguments()
	timestamp := time.Unix(int64(msg.Date), 0)

	var chunks []domain.PollChunk
	err := s.groupRepo.UpdateGroupAndPoll(ctx, msg.Chat.ID,
		func(group *domain.Group) (*domain.Group, *domain.Poll, error) {
			if group != nil {
				group.Name = msg.Chat.Title

				poll, err := group.CreatePoll(arg, timestamp, 10)
				if err != nil {
					return nil, nil, err
				}
				poll.MessageID = msg.MessageID

				chunks = poll.Chunks

				return group, poll, nil
			}

			return nil, nil, nil
		})

	if err == nil {
		for _, chunk := range chunks {
			poll := tgbotapi.NewPoll(msg.Chat.ID, chunk.Question, chunk.Options...)
			poll.IsAnonymous = false
			poll.ReplyToMessageID = msg.MessageID
			_, err = s.bot.Send(poll)
		}
	}
	if errors.Is(err, &domain.InvalidLeadtimeError{}) {
		reply := tgbotapi.NewMessage(msg.Chat.ID, err.Error())
		reply.ReplyToMessageID = msg.MessageID
		_, err = s.bot.Send(reply)
	}

	return err
}

func (s *LunchHandler) Handle() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()

		// For debug request body
		body := io.TeeReader(req.Body, os.Stdout)

		var update tgbotapi.Update
		err := json.NewDecoder(body).Decode(&update)
		if err == nil {
			err = s.process(req.Context(), &update)
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"status":"ERROR","message": "%v"}`, err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"status": "OK"}`)
	})
}
