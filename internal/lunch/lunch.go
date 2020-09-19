package lunch

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"giautm.dev/com/internal/bot"
	"giautm.dev/com/internal/lunch/domain"
)

func parseOptions(text string) (options []string, leadtime time.Duration) {
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

type BotPollSender struct {
	Bot      bot.ContextChat
	Question string
}

func (s *BotPollSender) SendMenu(ctx context.Context, items []domain.MenuItem) (string, error) {
	opts := make([]string, len(items))
	for idx, item := range items {
		opts[idx] = item.String()
	}

	return s.Bot.SendPoll(ctx, s.Question, opts)
}

func LunchQuestion(t time.Time, dur time.Duration) (string, time.Time) {
	closeDate := t.Add(dur)
	return fmt.Sprintf("%s Trưa nay ăn gì?\n\nChốt cơm lúc: %s",
		t.Format("2006-01-02"),
		closeDate.Format("15:04 02/01"),
	), closeDate
}

func (s *Server) OnCommand(bot bot.ContextChat, cmd string, args string) error {
	ctx := bot.Context()

	if cmd == "lunch" {
		opts, hours := parseOptions(args)

		menu := domain.NewMenuFlatPrice(opts, 30000)

		return s.repo.CreateMenuWithOrder(ctx, menu, func() (*domain.Order, error) {
			question, closeDate := LunchQuestion(time.Now(), hours)
			order, err := menu.Send(ctx, &BotPollSender{
				Bot:      bot,
				Question: question,
			})
			if err != nil {
				return nil, err
			}

			order.CloseDate = closeDate

			return order, nil
		})
	}
	return nil
}

func (s *Server) OnPollAnswer(ctx bot.Context, pollID string, options []int, userID int) error {
	// Retrict vote
	if len(options) == 0 {

	}
	fmt.Printf("OnPollAnswer: %s, %v, %d\n", pollID, options, userID)
	return nil
}
