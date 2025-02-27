package config

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	Cookies      []string `mapstructure:"cookies"`
	Password     string   `mapstructure:"password"`
	UserAgent    []string `mapstructure:"user_agent"`
	validAPIKeys map[string]bool
}

var (
	instance *Config
	once     sync.Once
)

// LoadConfig loads configuration from file
func LoadConfig() (*Config, error) {
	var err error
	once.Do(func() {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("./configs")
		viper.AddConfigPath(".")

		err = viper.ReadInConfig()
		if err != nil {
			return
		}

		instance = &Config{}
		err = viper.Unmarshal(instance)
		if err != nil {
			return
		}

		// Initialize valid API keys
		instance.validAPIKeys = make(map[string]bool)
		if instance.Password != "" {
			instance.validAPIKeys[instance.Password] = true
		}
	})

	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return instance, nil
}

// GetInstance returns the singleton instance of Config
func GetInstance() (*Config, error) {
	if instance == nil {
		return LoadConfig()
	}
	return instance, nil
}

// IsValidAPIKey checks if the API key is valid
func (c *Config) IsValidAPIKey(key string) bool {
	return c.validAPIKeys[key]
}
