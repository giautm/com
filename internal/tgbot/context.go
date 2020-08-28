package tgbot

import "context"

type chatContextKey struct{}

var key = chatContextKey{}

func WithChat(ctx context.Context, chatID int64) context.Context {
	return context.WithValue(ctx, key, chatID)
}

func ChatFromContext(ctx context.Context) int64 {
	val := ctx.Value(key)
	if val != nil {
		if chatID, ok := val.(int64); ok {
			return chatID
		}
	}
	return 0
}
