package main

import (
	"crypto-temka/internal/delivery"
	"crypto-temka/pkg/config"
	"crypto-temka/pkg/database"
	"crypto-temka/pkg/log"
)

func main() {
	logger, loggerInfoFile, loggerErrorFile := log.InitLogger()
	defer loggerInfoFile.Close()
	defer loggerErrorFile.Close()

	logger.Info("Logger Initialized")

	config.InitConfig()
	logger.Info("Config Initialized")

	db := database.MustGetDB()
	logger.Info("Database Initialized")

	delivery.Start(db, logger)
}
