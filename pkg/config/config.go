package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

const (
	DBName            = "DB_NAME"
	DBUser            = "DB_USER"
	DBPassword        = "DB_PASSWORD"
	DBPort            = "DB_PORT"
	DBHost            = "DB_HOST"
	Timeout           = "TIMEOUT"
	JWTExpire         = "JWT_EXPIRE"
	Secret            = "SECRET"
	SessionExpiration = "SESSION_EXPIRATION"
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