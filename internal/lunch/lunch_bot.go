package lunch

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"giautm.dev/com/internal/bot"
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

func (s *Server) OnCommand(ctx bot.ContextChat, cmd string, args string) error {
	if cmd == "lunch" {
		opts, hours := parseOptions(args)

		// Because Telegram Bot only support 10 minutes,
		// so We need close poll manually
		//
		// https://core.telegram.org/bots/api#sendpoll
		timestamp := time.Now()
		closeDate := timestamp.Add(hours)
		closeDateStr := closeDate.Format("15:04 02/01")
		question := fmt.Sprintf("* %s Trưa nay ăn gì?\n\nChốt cơm lúc: %s",
			timestamp.Format("2006-01-02"), closeDateStr)

		c := ctx.Context()
		pollID, err := ctx.SendPoll(c, question, opts)

		fmt.Printf("Poll ID: %s\n", pollID)
		return err
	}
	return nil
}

func (s *Server) OnPollAnswer(ctx bot.Context, pollID string, options []int, userID int) error {

	return nil
}
