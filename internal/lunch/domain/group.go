package domain

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type GroupRepo interface {
	UpdateGroup(
		ctx context.Context,
		chatId int64,
		updateFn func(*Group) (*Group, error),
	) error

	UpdateGroupAndPoll(
		ctx context.Context,
		chatID int64,
		updateFn func(*Group) (*Group, *Poll, error),
	) error
}

type Group struct {
	ID   string
	Name string

	Leadtime time.Duration
}

func (g *Group) parseOptions(text string) (options []string, leadtime time.Duration) {
	leadtime = g.Leadtime

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

func (g *Group) CreatePoll(input string, timestamp time.Time, maxOptions int) (*Poll, error) {
	options, hours := g.parseOptions(input)
	if hours <= 0 {
		return nil, &InvalidLeadtimeError{}
	}

	if len(options) > 0 {
		id, err := uuid.NewUUID()
		if err != nil {
			return nil, err
		}

		polls := chunkBy(options, maxOptions)
		totalPolls := len(polls)

		// Because Telegram Bot only support 10 minutes,
		// so We need close poll manually
		//
		// https://core.telegram.org/bots/api#sendpoll
		closeDate := timestamp.Add(hours)
		closeDateStr := closeDate.Format("15:04 02/01")

		chunks := make([]PollChunk, len(polls))
		for idx, pollOptions := range polls {
			question := fmt.Sprintf("%s Trưa nay ăn gì?", timestamp.Format("2006-01-02"))
			if totalPolls > 1 {
				question = fmt.Sprintf("[%d/%d] %s", idx+1, totalPolls, question)
			}
			question = fmt.Sprintf("%s\n\nChốt cơm lúc: %s", question, closeDateStr)

			chunks[idx] = PollChunk{
				Question: question,
				Options:  pollOptions,
			}
		}

		return &Poll{
			ID:        id,
			Input:     input,
			Chunks:    chunks,
			CreatedAt: timestamp,
			CloseDate: closeDate,
		}, nil
	}

	return nil, nil
}
