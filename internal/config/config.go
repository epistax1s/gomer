package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Report  ReportConfig
	Bot     BotConfig
	Redmine RedmineConfig
	Log     LogConfig
}

type ReportConfig struct {
	PublishCron      string
	NotificationCron string
}

type BotConfig struct {
	Username string
	Token    string
}

type RedmineConfig struct {
	BaseURL  string
	ApiKey   string
	Comments RedmineComments
}

type RedmineComments struct {
	Exclude []string
}

type LogConfig struct {
	Level  string
	Stdout bool
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	config := &Config{
		Report: ReportConfig{
			PublishCron:      getEnvRequired("REPORT_PUBLISH_CRON"),
			NotificationCron: getEnvRequired("REPORT_NOTIFICATION_CRON"),
		},
		Bot: BotConfig{
			Username: getEnvRequired("BOT_USERNAME"),
			Token:    getEnvRequired("BOT_TOKEN"),
		},
		Redmine: RedmineConfig{
			BaseURL: getEnvRequired("REDMINE_BASE_URL"),
			ApiKey:  getEnvRequired("REDMINE_API_KEY"),
			Comments: RedmineComments{
				Exclude: getEnvSliceRequired("REDMINE_COMMENTS_EXCLUDE", []string{
					"^митинг",
					"^ежедневный митинг",
					"^дейли",
					"^дэйли",
					"^daily",
				}),
			},
		},
		Log: LogConfig{
			Level:  getEnvRequired("LOG_LEVEL"),
			Stdout: getEnvRequiredBool("LOG_STDOUT"),
		},
	}

	return config, nil
}

// Helper functions for environment variables
func getEnvRequired(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("Required environment variable %s is not set", key))
	}
	return value
}

func getEnvRequiredBool(key string) bool {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("Required environment variable %s is not set", key))
	}
	if boolValue, err := strconv.ParseBool(value); err == nil {
		return boolValue
	} else {
		panic(fmt.Sprintf("Environment variable %s must be a valid boolean value (true/false, 1/0, yes/no), got: %s", key, value))
	}
}

func getEnvSliceRequired(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}
