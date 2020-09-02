package lunch

import (
	"context"
	"net/http"
)

func (s *Server) handleSchedule(ctx context.Context) http.HandlerFunc {
	bot := s.env.Bot()

	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), s.config.ScheduleTimeout)
		defer cancel()

		pollID := "123"
		bot.StopPoll(ctx, 0, pollID)
	}
}
