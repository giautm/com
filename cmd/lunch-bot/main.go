package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sethvargo/go-signalcontext"

	"giautm.dev/com/internal/lunch"
	"giautm.dev/com/internal/setup"
	"giautm.dev/com/pkg/logging"
	"giautm.dev/com/pkg/server"
)

func main() {
	ctx, done := signalcontext.OnInterrupt()

	logger := logging.NewLogger(true)
	ctx = logging.WithLogger(ctx, logger)

	err := realMain(ctx)
	done()

	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("successful shutdown")
}

func realMain(ctx context.Context) error {
	logger := logging.FromContext(ctx)

	var config lunch.Config
	env, err := setup.Setup(ctx, &config)
	if err != nil {
		return fmt.Errorf("setup.Setup: %w", err)
	}
	defer env.Close(ctx)

	handler, err := lunch.NewServer(&config, env)
	if err != nil {
		return fmt.Errorf("lunch.NewHandler: %w", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", handler.Routes(ctx))

	srv := server.New(config.Port)
	logger.Infof("listening on :%s", config.Port)

	return srv.ServeHTTPHandler(ctx, mux)
}
