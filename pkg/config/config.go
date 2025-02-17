package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	DBName     = "DB_NAME"
	DBUser     = "DB_USER"
	DBPassword = "DB_PASSWORD"
	DBPort     = "DB_PORT"
	DBHost     = "DB_HOST"
	Timeout    = "TIMEOUT"
	JWTExpire  = "JWT_EXPIRE"
	Secret     = "SECRET"

	ReferPercent = "REFER_PERCENT"

	OutcomeTickerMin = "OUTCOME_TICKER_MIN"
	OutcomeTickerMax = "OUTCOME_TICKER_MAX"

	OutcomeAmountMax = "OUTCOME_AMOUNT_MAX"

	OutcomeUserIDMin = "OUTCOME_USERID_MIN"
	OutcomeUserIDMax = "OUTCOME_USERID_MAX"

	OutcomesAmount = "OUTCOMES_AMOUNT"

	RateWorkerFrequency = "RATE_WORKER_FREQUENCY" // in seconds
)

func InitConfig() {
	envPath, _ := os.Getwd()

	envPath = filepath.Join(envPath, "..") // workdir is cmd
	envPath = filepath.Join(envPath, "/deploy")

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(envPath)

	// makes Viper check if environment variables match any of the existing keys (config, default or flags).
	//If matching env vars are found, they are loaded
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("Failed to init config. Error:%v", err.Error()))
	}
}
