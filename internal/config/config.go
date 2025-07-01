package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Report ReportConfig `json:"report"`
	Bot    BotConfig    `json:"bot"`
	Log    LogConfig    `json:"log"`
}

type ReportConfig struct {
	PublishCron      string `json:"publishCron"`
	NotificationCron string `json:"notificationCron"`
}

type BotConfig struct {
	Token string `json:"token"`
}

type LogConfig struct {
	Level  string `json:"level"`
	Stdout bool   `json:"stdout"`
}

func LoadConfig() (*Config, error) {
	bytes, err := os.ReadFile("./config/config.json")
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
