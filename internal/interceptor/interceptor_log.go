package interceptor

import (
	"github.com/epistax1s/gomer/internal/log"
	"github.com/epistax1s/gomer/internal/server"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type LogInterceptor struct {
	BaseInterceptor
}

func (i *LogInterceptor) Handle(server *server.Server, update *tgbotapi.Update) {
	log.Debug("[LogInterceptor] - Update", update)
	if update.CallbackQuery != nil {
		log.Debug("[LogInterceptor] - CallbackQueryData", update.CallbackQuery.Data)
	}

	i.Next(server, update)
}
