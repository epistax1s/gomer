package state

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/epistax1s/gomer/internal/i18n"
	"github.com/epistax1s/gomer/internal/log"
	"github.com/epistax1s/gomer/internal/server"
)

type CommitState struct {
	data *StateContext
}

func NewCommitState(data *StateContext) State {
	return &CommitState{
		data: data,
	}
}

func (state *CommitState) Init(server *server.Server, update *tgbotapi.Update) {
	chatID := update.FromChat().ID

	commit, err := server.CommitService.FindCommitByUserIdAndDate(chatID, state.data.CommitDate)
	if err != nil {
		log.Error(
			"error when searching for a commit",
			"chatID", chatID, "state", commit, "step", "Init", "err", err)

		server.Gomer.SendMessage(chatID, "Что-то пошло не так.")

		StateMachine.
			Set(Idle, chatID, &StateContext{}).
			Handle(server, update)

		return
	}

	if commit != nil {
		log.Info(
			"attempt to create a second commit on the same date",
			"chatID", chatID, "state", commit, "step", "Init", "date", state.data.CommitDate)

		msg := fmt.Sprintf(i18n.Localize("commitAlreadyCreated"), state.data.CommitDate, commit.Payload)

		server.Gomer.SendMessage(chatID, msg)

		StateMachine.
			Set(Idle, chatID, &StateContext{}).
			Init(server, update)

		return
	}

	msg := fmt.Sprintf(i18n.Localize("commitCreatePromt"), state.data.CommitDate)

	server.Gomer.SendMessage(chatID, msg)
}

func (state *CommitState) Handle(server *server.Server, update *tgbotapi.Update) {
	chatID := update.FromChat().ID

	// validate input
	if update.Message == nil {
		log.Error(
			"message is nil",
			"chatID", chatID, "state", Commit)

		state.Init(server, update)
		return
	}

	payload := update.Message.Text

	// create commit
	err := server.CommitService.CreateCommit(chatID, payload, state.data.CommitDate)
	if err != nil {
		server.Gomer.SendMessage(chatID, i18n.Localize("commitCreateError"))
		return
	}

	// send success response
	msg := fmt.Sprintf(i18n.Localize("commitCreateSuccess"), state.data.CommitDate, payload)

	server.Gomer.SendMessage(chatID, msg)

	// set base state
	StateMachine.
		Set(Idle, chatID, &StateContext{}).
		Init(server, update)
}
