package state

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/epistax1s/gomer/internal/i18n"
	"github.com/epistax1s/gomer/internal/log"
	"github.com/epistax1s/gomer/internal/model"
	"github.com/epistax1s/gomer/internal/server"
)

type TrackNameState struct {
	data *StateContext
}

func NewTrackNameState(data *StateContext) State {
	return &TrackNameState{
		data: data,
	}
}

func (state *TrackNameState) Init(server *server.Server, update *tgbotapi.Update) {
	server.Gomer.SendMessage(update.FromChat().ID, i18n.Localize("enterNamePromt"))
}

func (state *TrackNameState) Handle(server *server.Server, update *tgbotapi.Update) {
	chatID := update.FromChat().ID
	chatUsername := update.FromChat().UserName

	if update.Message == nil {
		log.Error(
			"Message is nil",
			"chatID", chatID, "state", TrackName, "step", "Handle")

		state.Init(server, update)
		return
	}

	messageText := update.Message.Text

	name := strings.TrimSpace(messageText)
	if !isValidName(name) {
		log.Info(
			"Name has an invalid value",
			"chatID", chatID, "state", TrackName, "step", "Handle", "name", name)

		server.Gomer.SendMessage(chatID, i18n.Localize("enterNameNotValid"))

		state.Init(server, update)
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

	err := server.UserService.TrackUser(user)
	if err != nil {
		log.Info(
			"error registering a new user",
			"chatID", chatID, "state", TrackName, "step", "Handle", "name", name)

		server.Gomer.SendMessage(chatID, i18n.Localize("oops"))

		StateMachine.
			Set(Idle, chatID, &StateContext{}).
			Init(server, update)
	}

	server.Gomer.SendMessage(chatID, i18n.Localize("trackStarted"))

	log.Info(
		"the user has successfully registered",
		"chatID", chatID, "state", TrackName, "step", "Handle", "name", name)

	StateMachine.
		Set(Idle, chatID, &StateContext{}).
		Init(server, update)
}

func isValidName(name string) bool {
	parts := strings.Split(name, " ")
	return len(parts) == 2
}
