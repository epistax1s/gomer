package state

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/epistax1s/gomer/internal/i18n"
	"github.com/epistax1s/gomer/internal/server"
)

type IdleState struct {
}

func NewIdleState(data *StateContext) State {
	return &IdleState{}
}

func (state *IdleState) Init(server *server.Server, update *tgbotapi.Update) {
	helpHandler := IdleCmdHandlers["help"]
	helpHandler(server, update)
}

func (state *IdleState) Handle(server *server.Server, update *tgbotapi.Update) {
	chatID := update.FromChat().ID
	cmd := update.Message.Command()

	cmdHandler, exits := IdleCmdHandlers[cmd]
	if exits {
		cmdHandler(server, update)
	} else {
		server.Gomer.SendMessage(chatID, i18n.Localize("unsupportedCommand"))
		state.Init(server, update)
	}
}
