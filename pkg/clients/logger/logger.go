package logger

import (
	"os"

	"github.com/ContextLogic/pkg/autobots/config"
	"github.com/sirupsen/logrus"
)

// New creates a new logger client
func New(config *config.LoggerConfig) (*logrus.Logger, error) {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	level, _ := logrus.ParseLevel(config.Level)
	logger.SetLevel(level)
	logger.SetOutput(os.Stdout)
	return logger, nil
}
