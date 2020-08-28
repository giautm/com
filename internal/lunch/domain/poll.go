package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type PollRepo interface {
	Save(ctx *context.Context, poll *Poll) error
}

type Poll struct {
	ID        uuid.UUID
	MessageID int
	Input     string
	Chunks    []PollChunk
	CreatedAt time.Time
	CloseDate time.Time
}

type PollChunk struct {
	Question string
	Options  []string
	IsSend   bool
	IsClosed bool

	PollID string
}
