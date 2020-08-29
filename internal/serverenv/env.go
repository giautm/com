// Package serverenv defines common parameters for the sever environment.
package serverenv

import (
	"context"

	"giautm.dev/com/internal/bot"
)

// ServerEnv represents latent environment configuration for servers in this application.
type ServerEnv struct {
	bot bot.Bot
}

// Option defines function types to modify the ServerEnv on creation.
type Option func(*ServerEnv) *ServerEnv

// New creates a new ServerEnv with the requested options.
func New(ctx context.Context, opts ...Option) *ServerEnv {
	env := &ServerEnv{}
	for _, f := range opts {
		env = f(env)
	}

	return env
}

// // WithDatabase attached a database to the environment.
// func WithDatabase(db *database.DB) Option {
// 	return func(s *ServerEnv) *ServerEnv {
// 		s.database = db
// 		return s
// 	}
// }

func WithBot(p bot.Bot) Option {
	return func(s *ServerEnv) *ServerEnv {
		s.bot = p
		return s
	}
}

func (s *ServerEnv) Bot() bot.Bot {
	return s.bot
}

// Close shuts down the server env, closing database connections, etc.
func (s *ServerEnv) Close(ctx context.Context) error {
	if s == nil {
		return nil
	}

	// if s.database != nil {
	// 	s.database.Close(ctx)
	// }

	return nil
}
