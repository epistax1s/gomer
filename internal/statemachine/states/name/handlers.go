package name

import (
	"strings"

	"github.com/epistax1s/gomer/internal/i18n"
	"github.com/epistax1s/gomer/internal/log"
	"github.com/epistax1s/gomer/internal/model"

	. "github.com/epistax1s/gomer/internal/statemachine/core"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (state *TrackNameState) Init(update *tgbotapi.Update) {
	gomer := state.server.Gomer
	gomer.SendMessage(update.FromChat().ID, i18n.Localize("enterNamePromt"))
}

func (state *TrackNameState) Handle(update *tgbotapi.Update) {
	gomer := state.server.Gomer
	userService := state.server.UserService

	chatID := update.FromChat().ID
	chatUsername := update.FromChat().UserName

	if update.Message == nil {
		log.Error(
			"Message is nil",
			"chatID", chatID, "state", TrackName, "step", "Handle")

		state.Init(update)
		return
	}

	messageText := update.Message.Text

	name := strings.TrimSpace(messageText)
	if !isValidName(name) {
		log.Info(
			"Name has an invalid value",
			"chatID", chatID, "state", TrackName, "step", "Handle", "name", name)

		gomer.SendMessage(chatID, i18n.Localize("enterNameNotValid"))

		state.Init(update)
		return
	}

	user := &model.User{
		ChatID:       chatID,
		DepartmentId: state.data.Department.ID,
		Order:        0,
		Name:         messageText,
		Username:     chatUsername,
		Role:         model.UserRoleUser,
		Status:       model.UserStatusActive,
	}

	err := userService.TrackUser(user)
	if err != nil {
		log.Info(
			"error registering a new user",
			"chatID", chatID, "state", TrackName, "step", "Handle", "name", name)

		gomer.SendMessage(chatID, i18n.Localize("oops"))

		state.stateMachine.
			Set(Idle, chatID, &StateContext{}).
			Init(update)
	}

	gomer.SendMessage(chatID, i18n.Localize("trackStarted"))

	log.Info(
		"the user has successfully registered",
		"chatID", chatID, "state", TrackName, "step", "Handle", "name", name)

	state.stateMachine.
		Set(Idle, chatID, &StateContext{}).
		Init(update)
}

func isValidName(name string) bool {
	parts := strings.Split(name, " ")
	return len(parts) == 2
}
