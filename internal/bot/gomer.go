package gomer

import (
	"github.com/epistax1s/gomer/internal/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Gomer struct {
	tgbotapi.BotAPI
}

func InitTelegramBot(botConfig *config.BotConfig) (*Gomer, error) {
	botApi, err := tgbotapi.NewBotAPI(botConfig.Token)
	if err != nil {
		return nil, err
	}
	return &Gomer{BotAPI: *botApi}, nil
}
