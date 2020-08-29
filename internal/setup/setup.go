// Package setup provides common logic for configuring the various services.
package setup

import (
	"context"
	"fmt"

	"github.com/sethvargo/go-envconfig"

	"giautm.dev/com/internal/bot"
	"giautm.dev/com/internal/serverenv"
	"giautm.dev/com/pkg/logging"
)

// BotConfigProvider signals that the config provided knows how to
// configure bot.
type BotConfigProvider interface {
	BotConfig() *bot.Config
}

// // DatabaseConfigProvider ensures that the environment config can provide a DB config.
// // All binaries in this application connect to the database via the same method.
// type DatabaseConfigProvider interface {
// 	DatabaseConfig() *database.Config
// }

// Setup runs common initialization code for all servers. See SetupWith.
func Setup(ctx context.Context, config interface{}) (*serverenv.ServerEnv, error) {
	return SetupWith(ctx, config, envconfig.OsLookuper())
}

// SetupWith processes the given configuration using envconfig. It is
// responsible for establishing database connections, resolving secrets, and
// accessing app configs. The provided interface must implement the various
// interfaces.
func SetupWith(ctx context.Context, config interface{}, l envconfig.Lookuper) (*serverenv.ServerEnv, error) {
	logger := logging.FromContext(ctx)

	// Build a list of mutators. This list will grow as we initialize more of the
	// configuration, such as the secret manager.
	var mutatorFuncs []envconfig.MutatorFunc

	// Build a list of options to pass to the server env.
	var serverEnvOpts []serverenv.Option

	// Process first round of environment variables.
	if err := envconfig.ProcessWith(ctx, config, l, mutatorFuncs...); err != nil {
		return nil, fmt.Errorf("error loading environment variables: %w", err)
	}
	logger.Infow("provided", "config", config)

	// Configure blob storage.
	if provider, ok := config.(BotConfigProvider); ok {
		logger.Info("configuring bot")

		botCfg := provider.BotConfig()
		bot, err := bot.BotFor(ctx, botCfg)
		if err != nil {
			return nil, fmt.Errorf("unable to connect to storage system: %v", err)
		}

		// Update serverEnv setup.
		serverEnvOpts = append(serverEnvOpts, serverenv.WithBot(bot))

		logger.Infow("bot", "config", botCfg)
	}

	return serverenv.New(ctx, serverEnvOpts...), nil
}
