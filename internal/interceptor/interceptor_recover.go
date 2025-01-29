package interceptor

import (
	"github.com/epistax1s/gomer/internal/log"
	"github.com/epistax1s/gomer/internal/server"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type RecoverInterceptor struct {
	BaseInterceptor
}

func (i *RecoverInterceptor) Handle(server *server.Server, update *tgbotapi.Update) {
	defer func() {
		if r := recover(); r != nil {
			log.Error("Panic is caught: ", "panic", r)
		}
	}()

	i.Next(server, update)
}
