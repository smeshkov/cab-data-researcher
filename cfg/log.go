package cfg

import (
	"fmt"

	"go.uber.org/zap"
)

// LogSetup configures logger.
func LogSetup(env, version string) error {
	l, err := newLog(env)
	if err != nil {
		return fmt.Errorf("failed to setup logger: %w", err)
	}

	defer LogSync()

	if env == "" {
		env = "local"
	}
	zap.ReplaceGlobals(l.With(
		zap.String("env", env),
		zap.String("version", version),
	))
	zap.L().Info("logger is ready")

	return nil
}

// LogSync flushes any buffered log entries.
func LogSync() {
	_ = zap.L().Sync()
}

func newLog(env string) (*zap.Logger, error) {
	if env != "prod" {
		return zap.NewDevelopment()
	}
	return zap.NewProduction()
}
