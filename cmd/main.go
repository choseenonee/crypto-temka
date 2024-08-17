package main

import (
	"crypto-temka/internal/delivery"
	"crypto-temka/internal/service"
	"crypto-temka/pkg/config"
	"crypto-temka/pkg/database"
	"crypto-temka/pkg/log"
	"fmt"
	"os"
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

	metricsSetFile, err := os.OpenFile("static/metrics.json", os.O_RDWR, 0660)
	if err != nil {
		panic(fmt.Sprintf("Error opening static/metrics.json file, err: %v", err.Error()))
	}
	defer metricsSetFile.Close()

	service.InitUserRateWorker(db, logger)
	delivery.Start(db, logger, metricsSetFile)
}
