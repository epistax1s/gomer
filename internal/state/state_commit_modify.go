package state

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/epistax1s/gomer/internal/i18n"
	"github.com/epistax1s/gomer/internal/log"
	"github.com/epistax1s/gomer/internal/server"
)

type CommitModifyState struct {
	data *StateContext
}

func NewCommitModifyState(data *StateContext) State {
	return &CommitModifyState{
		data: data,
	}
}

func (state *CommitModifyState) Init(server *server.Server, update *tgbotapi.Update) {
	chatID := update.FromChat().ID

	commit, err := server.CommitService.FindCommitByUserIdAndDate(chatID, state.data.CommitDate)
	if err != nil {
		log.Error(
			"error when searching for a commit",
			"chatID", chatID, "state", CommitModify, "step", "Init", "err", err)

		server.Gomer.SendMessage(chatID, i18n.Localize("oops"))

		StateMachine.
			Set(Idle, chatID, &StateContext{}).
			Handle(server, update)

		return
	}

	if commit == nil {
		msg := fmt.Sprintf(i18n.Localize("commitNotFound"), state.data.CommitDate)

		server.Gomer.SendMessage(chatID, msg)

		StateMachine.
			Set(Idle, chatID, &StateContext{}).
			Handle(server, update)

		return
	}

	state.data.Commit = commit

	msg := fmt.Sprintf(i18n.Localize("commitModifyPromt"), commit.Payload, commit.Date)

	server.Gomer.SendMessage(chatID, msg)
}

func (state *CommitModifyState) Handle(server *server.Server, update *tgbotapi.Update) {
	chatID := update.FromChat().ID

	// validate input
	if update.Message == nil {
		log.Error(
			"message is nil",
			"chatID", chatID, "state", Commit)

		state.Init(server, update)
		return
	}

	// extract message
	payload := update.Message.Text

	// update commit
	commit, err := server.CommitService.UpdateCommit(state.data.Commit.ID, payload)
	if err != nil {
		log.Error(
			"error updating commit",
			"chatID", chatID, "err", err)

		server.Gomer.SendMessage(chatID, i18n.Localize("commitModifyError"))

		StateMachine.
			Set(Idle, chatID, &StateContext{}).
			Init(server, update)

		return
	}

	// notify about success
	msg := fmt.Sprintf(i18n.Localize("commitModifySuccess"), commit.Date, commit.Payload)

	server.Gomer.SendMessage(chatID, msg)

	// set starting state
	StateMachine.
		Set(Idle, chatID, &StateContext{}).
		Init(server, update)
}
