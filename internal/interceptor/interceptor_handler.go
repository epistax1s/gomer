package interceptor

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/epistax1s/gomer/internal/i18n"
	"github.com/epistax1s/gomer/internal/log"
	"github.com/epistax1s/gomer/internal/server"
	"github.com/epistax1s/gomer/internal/state"
)

type HandlerInterceptor struct {
	BaseInterceptor
}

func (i *HandlerInterceptor) Handle(server *server.Server, update *tgbotapi.Update) {
	if update == nil {
		return
	}

	if update.FromChat().IsPrivate() {
		handleFromPrivate(server, update)
	} else if update.FromChat().IsGroup() || update.FromChat().IsSuperGroup() {
		handleFromGroup(server, update)
	} else {
		log.Error("This bot can only be used in private conversations and in groups")
	}
}

func handleFromPrivate(server *server.Server, update *tgbotapi.Update) {
	chatID := update.FromChat().ID

	state.StateMachine.
		Get(chatID).
		Handle(server, update)
}

func handleFromGroup(server *server.Server, update *tgbotapi.Update) {
	chatID := update.FromChat().ID
	title := update.FromChat().Title

	cmd := update.Message.Command()

	if cmd == "link" {
		err := server.GroupService.LinkGroup(chatID, title)
		if err == nil {
			log.Info("a new group has been assigned to the bot",
				"chatID", chatID, "title", title)

			server.Gomer.SendMessage(chatID, i18n.Localize("groupSuccessfullyLinked"))
		}
	}
}
