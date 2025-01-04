package interceptor

import (
	"github.com/epistax1s/gomer/internal/log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type LogInterceptor struct {
	BaseInterceptor
}

func (i *LogInterceptor) Handle(update *tgbotapi.Update) {
	log.Debug("[LogInterceptor] - Update", update)
	if update.CallbackQuery != nil {
		log.Debug("[LogInterceptor] - CallbackQueryData", update.CallbackQuery.Data)
	}

	i.Next(update)
}
