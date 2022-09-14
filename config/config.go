package config

import (
	"log"

	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variables.
type Config struct {
	NASAAPIKey         string `mapstructure:"API_KEY"`
	ConcurrentRequests string `mapstructure:"CONCURRENT_REQUESTS"`
	ServerPort         string `mapstructure:"PORT"`
	DBDriver           string `mapstructure:"DB_DRIVER"`
	DBSource           string `mapstructure:"DB_SOURCE"`
	MigrationURL       string `mapstructure:"MIGRATION_URL"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		log.Printf("failed to read config variables by viper: %v", err)
		return
	}

	err = viper.Unmarshal(&config)

	return
}
