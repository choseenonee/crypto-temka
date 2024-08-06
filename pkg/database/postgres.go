package database

import (
	"crypto-temka/pkg/config"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"

	_ "github.com/lib/pq"
)

func MustGetDB() *sqlx.DB {
	connString := fmt.Sprintf(
		"user=%v password=%v host=%v port=%v dbname=%v sslmode=disable",
		viper.GetString(config.DBUser),
		viper.GetString(config.DBPassword),
		viper.GetString(config.DBHost),
		viper.GetInt(config.DBPort),
		viper.GetString(config.DBName),
	)
	db, err := sqlx.Connect("postgres", connString)
	if err != nil {
		panic(fmt.Sprintf("Error while connecting to DB. Error: %v", err.Error()))
	}

	return db
}
