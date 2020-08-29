package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"giautm.dev/com/pkg/logging"
	srvcloud "gocloud.dev/server"
	"gocloud.dev/server/requestlog"
)

type Server interface {
	ServeHTTPHandler(ctx context.Context, handler http.Handler) error
}

type server struct {
	port string
}

func New(port string) Server {
	return &server{
		port: port,
	}
}

func (s *server) ServeHTTPHandler(ctx context.Context, handler http.Handler) error {
	logger := logging.FromContext(ctx)

	// Create a logger, and assign it to the RequestLogger field of a
	// server.Options struct.
	srvOptions := &srvcloud.Options{
		RequestLogger: requestlog.NewNCSALogger(os.Stdout, func(error) {}),
	}
	srv := srvcloud.New(handler, srvOptions)

	// Spawn a goroutine that listens for context closure. When the context is
	// closed, the server is stopped.
	errCh := make(chan error, 1)
	go func() {
		// Wait for CTRL+C
		<-ctx.Done()

		logger.Debugf("server.Serve: context closed")
		shutdownCtx, done := context.WithTimeout(context.Background(), 5*time.Second)
		defer done()

		logger.Debugf("server.Serve: shutting down")
		if err := srv.Shutdown(shutdownCtx); err != nil {
			select {
			case errCh <- err:
			default:
			}
		}
	}()

	// Run the server. This will block until the provided context is closed.
	if err := srv.ListenAndServe(":" + s.port); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to serve: %w", err)
	}

	logger.Debugf("server.Serve: serving stopped")

	// Return any errors that happened during shutdown.
	select {
	case err := <-errCh:
		return fmt.Errorf("failed to shutdown: %w", err)
	default:
		return nil
	}
}
