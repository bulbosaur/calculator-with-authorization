package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

// Init считывает переменные окружения
func Init() {

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetDefault("server.HTTP_HOST", "localhost")
	viper.SetDefault("server.HTTP_PORT", "8080")
	viper.SetDefault("server.GRPC_HOST", "localhost")
	viper.SetDefault("server.GRPC_PORT", "50051")

	viper.SetDefault("duration.TIME_ADDITION_MS", 100)
	viper.SetDefault("duration.TIME_SUBTRACTION_MS", 100)
	viper.SetDefault("duration.TIME_MULTIPLICATIONS_MS", 100)
	viper.SetDefault("duration.TIME_DIVISIONS_MS", 100)
	viper.SetDefault("DATABASE_PATH", "./db/calc.db")
	viper.SetDefault("worker.COMPUTING_POWER", 5)

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath("./config")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading .env file: %v", err)
	}

	logConfig()
}

func logConfig() {
	log.Printf(
		"Configuration: HTTP_HOST=%s, HTTP_PORT=%s, GRPC_HOST=%s, GRPC_PORT=%s, TIME_ADDITION_MS=%d, TIME_SUBTRACTION_MS=%d, TIME_MULTIPLICATIONS_MS=%d, TIME_DIVISIONS_MS=%d, DATABASE_PATH=%s",
		viper.GetString("server.HTTP_HOST"),
		viper.GetString("server.HTTP_PORT"),
		viper.GetString("server.GRPC_HOST"),
		viper.GetString("server.GRPC_PORT"),
		viper.GetInt("duration.TIME_ADDITION_MS"),
		viper.GetInt("duration.TIME_SUBTRACTION_MS"),
		viper.GetInt("duration.TIME_MULTIPLICATIONS_MS"),
		viper.GetInt("duration.TIME_DIVISIONS_MS"),
		viper.GetString("DATABASE_PATH"),
	)
}
