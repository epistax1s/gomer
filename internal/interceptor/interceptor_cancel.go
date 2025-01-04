package interceptor

import (
	"github.com/epistax1s/gomer/internal/log"
	"github.com/epistax1s/gomer/internal/server"
	"github.com/epistax1s/gomer/internal/statemachine/core"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CancelInterceptor struct {
	BaseInterceptor
	Server       *server.Server
	StateMachine *core.StateMachine
}

func (i *CancelInterceptor) Handle(update *tgbotapi.Update) {
	if update != nil && update.FromChat().IsPrivate() && update.Message != nil {
		cmd := update.Message.Command()
		if cmd == "cancel" {
			chatID := update.FromChat().ID

			i.StateMachine.
				Set(core.Idle, chatID, &core.StateContext{}).
				Init(update)

			log.Info(
				"the status has been reset for the user",
				"chatID", chatID)

			return
		}
	}

	i.Next(update)
}
