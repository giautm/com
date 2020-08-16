package domain

import "context"

type Group struct {
	Name string
}

type GroupRepo interface {
	GetOrCreateGroup(ctx context.Context, chatID int64) (*Group, error)
}
