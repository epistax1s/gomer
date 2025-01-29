package interceptor

import (
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/epistax1s/gomer/internal/log"
	"github.com/epistax1s/gomer/internal/server"
	"github.com/epistax1s/gomer/internal/state"
)

type HandlerInterceptor struct {
	BaseInterceptor
}

func (i *HandlerInterceptor) Handle(server *server.Server, update *tgbotapi.Update) {
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
	chat := update.FromChat()
	chatID := chat.ID

	cmd := update.Message.Command()

	if cmd == "link" {
		err := server.GroupService.LinkGroup(chatID, chat.Title)
		if err == nil {
			log.Info("ok we added this group to report target, msg.From.ID=", chatID)
			rmsg := tgbotapi.NewMessage(chatID, "ok we added this group to report target, msg.From.ID = "+strconv.FormatInt(chatID, 10))
			server.Gomer.Send(rmsg)
		}
	}
}