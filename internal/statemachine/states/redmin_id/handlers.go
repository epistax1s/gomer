package redminid

import (
	"strconv"
	"strings"

	"github.com/epistax1s/gomer/internal/i18n"
	"github.com/epistax1s/gomer/internal/log"

	. "github.com/epistax1s/gomer/internal/statemachine/core"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (state *RedmineIDState) Init(update *tgbotapi.Update) {
	gomer := state.server.Gomer
	chatID := update.FromChat().ID

	gomer.SendMessageHtml(chatID, i18n.Localize("redmineIDRequest_html"))
}

func (state *RedmineIDState) Handle(update *tgbotapi.Update) {
	gomer := state.server.Gomer
	userService := state.server.UserService
	chatID := update.FromChat().ID

	if update.Message == nil {
		log.Error(
			"Message is nil",
			"chatID", chatID, "state", "RedmineID", "step", "Handle")
		state.Init(update)
		return
	}

	msg := strings.TrimSpace(update.Message.Text)
	redmineID, isValid := parseRedmineId(msg)
	if !isValid {
		log.Warn(
			"RedmineID has an invalid value",
			"chatID", chatID, "state", "RedmineID", "msg", msg)

		gomer.SendMessage(chatID, i18n.Localize("redmineIDNotValid"))
		state.Init(update)
		return
	}

	user, _ := userService.FindByChatID(chatID)
	user.RedmineID = &redmineID
	userService.Save(user)

	gomer.SendMessage(chatID, i18n.Localize("redmineIDSavedSuccess"))

	state.stateMachine.
		Set(Config, chatID, &StateContext{}).
		Init(update)
}

func parseRedmineId(id string) (int64, bool) {
	redmineID, err := strconv.ParseInt(id, 10, 64)

	if err != nil || redmineID < 1 {
		return 0, false
	}

	return redmineID, true
}
