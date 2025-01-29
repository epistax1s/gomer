package interceptor

import (
	"github.com/epistax1s/gomer/internal/log"
	"github.com/epistax1s/gomer/internal/server"
	"github.com/epistax1s/gomer/internal/state"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CancelInterceptor struct {
	BaseInterceptor
}

func (i *CancelInterceptor) Handle(server *server.Server, update *tgbotapi.Update) {
	if update != nil && update.FromChat().IsPrivate() && update.Message != nil  {
		cmd := update.Message.Command()
		if cmd == "cancel" {
			chatID := update.FromChat().ID
			
			state.StateMachine.
				Set(state.Idle, chatID, &state.StateContext{}).
				Init(server, update)

			log.Info(
				"the status has been reset for the user", 
				"chatID", chatID)

			return
		}
	} 

	i.Next(server, update);
}
