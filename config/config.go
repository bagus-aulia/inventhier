package config

import (
	"log"

	"github.com/spf13/viper"
)

// Config holds all configuration properties loaded from environment variables or .env.
type Config struct {
	AppEnv     string `mapstructure:"APP_ENV" default:"development"`
	ServerPort string `mapstructure:"SERVER_PORT" default:"8080"`

	SQLHost string `mapstructure:"SQL_HOST" default:"localhost"`
	SQLPort string `mapstructure:"SQL_PORT" default:"5432"`
	SQLUser string `mapstructure:"SQL_USER" default:"root"`
	SQLPass string `mapstructure:"SQL_PASS" default:""`
	SQLName string `mapstructure:"SQL_DB_NAME" default:"inventhier_db"`

	RedisHost string `mapstructure:"REDIS_HOST" default:"localhost"`
	RedisPort string `mapstructure:"REDIS_PORT" default:"6379"`

	MongoURI    string `mapstructure:"MONGO_URI" default:"mongodb://root:password@localhost:27017"`
	MongoDBName string `mapstructure:"MONGO_DB_NAME" default:"inventhier_db"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig() (*Config, error) {
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	// If a .env file exists, read it
	if err := viper.ReadInConfig(); err != nil {
		log.Println("No .env file found, using system environment variables instead")
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
