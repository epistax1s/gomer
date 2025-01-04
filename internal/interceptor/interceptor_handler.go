package interceptor

import (
	"github.com/epistax1s/gomer/internal/i18n"
	"github.com/epistax1s/gomer/internal/log"
	"github.com/epistax1s/gomer/internal/server"
	"github.com/epistax1s/gomer/internal/statemachine/core"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type HandlerInterceptor struct {
	Server       *server.Server
	StateMachine *core.StateMachine
	BaseInterceptor
}

func (i *HandlerInterceptor) Handle(update *tgbotapi.Update) {
	if update == nil {
		return
	}

	if update.FromChat().IsPrivate() {
		i.handleFromPrivate(update)
	} else if update.FromChat().IsGroup() || update.FromChat().IsSuperGroup() {
		i.handleFromGroup(update)
	} else {
		log.Error("This bot can only be used in private conversations and in groups")
	}
}

func (i *HandlerInterceptor) handleFromPrivate(update *tgbotapi.Update) {
	chatID := update.FromChat().ID

	i.StateMachine.
		Get(chatID).
		Handle(update)
}

func (i *HandlerInterceptor) handleFromGroup(update *tgbotapi.Update) {
	gomer := i.Server.Gomer
	groupService := i.Server.GroupService

	chatID := update.FromChat().ID
	title := update.FromChat().Title

	cmd := update.Message.Command()

	if cmd == "link" {
		if err := groupService.LinkGroup(chatID, title); err == nil {
			gomer.SendMessage(chatID, i18n.Localize("groupSuccessfullyLinked"))
		} else {
			gomer.SendMessage(chatID, i18n.Localize("oops"))
		}
	}
}
