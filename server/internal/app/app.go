package app

import (
	"messenger/internal/config"
	"messenger/internal/logger/sl"
)

func Run(cfg *config.Config) {
	log := sl.InitLogger(cfg.Env)

	log.Info("Logger is enabled")
	log.Debug("Debug is enabled")
}
