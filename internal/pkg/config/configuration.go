package config

import (
	"log"

	"github.com/spf13/viper"
)

type DatabaseConfig struct {
	Driver       string
	Dbname       string
	Username     string
	Password     string
	Host         string
	Port         string
	MaxLifetime  int
	MaxOpenConns int
	MaxIdleConns int
}

type ServerConfig struct {
	Port   string
	Secret string
	Mode   string
}

type Configuration struct {
	Server   ServerConfig
	Database DatabaseConfig
}

var Config *Configuration

// GetConfig get configuration data
func GetConfig() *Configuration {
	return Config
}

// Setup initialize configuration
func Setup(configPath string) {
	var configuration *Configuration
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(&configuration)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	Config = configuration
}
