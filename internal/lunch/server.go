package lunch

import (
	"context"
	"fmt"
	"net/http"

	"giautm.dev/com/internal/lunch/domain"
	"giautm.dev/com/internal/serverenv"
)

// Server hosts end points to manage export batches.
type Server struct {
	config *Config
	env    *serverenv.ServerEnv

	repo domain.MenuRepo
}

// NewServer makes a Server.
func NewServer(config *Config, env *serverenv.ServerEnv) (*Server, error) {
	// Validate config.
	if env.Bot() == nil {
		return nil, fmt.Errorf("lunch.NewServer requires Bot present in the ServerEnv")
	}

	return &Server{
		config: config,
		env:    env,
	}, nil
}

// Routes defines and returns the routes for this server.
func (s *Server) Routes(ctx context.Context) *http.ServeMux {
	mux := http.NewServeMux()

	bot := s.env.Bot()
	mux.HandleFunc(bot.WebhookPath(), bot.UpdatesHandler(ctx, s))
	mux.HandleFunc("/run-schedule", s.handleSchedule(ctx))

	return mux
}
