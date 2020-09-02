package bot

import (
	"context"
)

type BotRepo interface {
	SavePoll(ctx context.Context, pollID string, MessageID int) error
	FetchMessageID(ctx context.Context, pollID string) (int, error)
}

type memBotRepo struct {
	polls map[string]int
}

func NewMemoryBotRepo() BotRepo {
	return &memBotRepo{
		polls: make(map[string]int),
	}
}

func (m *memBotRepo) SavePoll(_ context.Context, pollID string, msgID int) error {
	m.polls[pollID] = msgID
	return nil
}

func (m memBotRepo) FetchMessageID(_ context.Context, pollID string) (int, error) {
	return m.polls[pollID], nil
}
