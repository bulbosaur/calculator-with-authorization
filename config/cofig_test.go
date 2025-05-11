package config_test

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"github.com/bulbosaur/calculator-with-authorization/config"
)

func TestDefaultConfig(t *testing.T) {
	os.Clearenv()
	config.Init()

	assert.Equal(t, "localhost", viper.GetString("server.HTTP_HOST"))
	assert.Equal(t, "8080", viper.GetString("server.HTTP_PORT"))
	assert.Equal(t, "localhost", viper.GetString("server.GRPC_HOST"))
	assert.Equal(t, "50051", viper.GetString("server.GRPC_PORT"))

	assert.Equal(t, 100, viper.GetInt("duration.TIME_ADDITION_MS"))
	assert.Equal(t, 100, viper.GetInt("duration.TIME_SUBTRACTION_MS"))
	assert.Equal(t, 100, viper.GetInt("duration.TIME_MULTIPLICATIONS_MS"))
	assert.Equal(t, 100, viper.GetInt("duration.TIME_DIVISIONS_MS"))

	assert.Equal(t, "./db/calc.db", viper.GetString("DATABASE_PATH"))
	assert.Equal(t, 5, viper.GetInt("worker.COMPUTING_POWER"))

	assert.Equal(t, "your_secret_key_here", viper.GetString("jwt.secret_key"))
	assert.Equal(t, 24, viper.GetInt("jwt.token_duration"))
}
