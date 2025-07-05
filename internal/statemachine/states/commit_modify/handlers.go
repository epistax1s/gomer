package commit_modify

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/epistax1s/gomer/internal/i18n"
	"github.com/epistax1s/gomer/internal/log"

	. "github.com/epistax1s/gomer/internal/statemachine/core"
)

func (state *CommitModifyState) Init(update *tgbotapi.Update) {
	gomer := state.server.Gomer
	commitService := state.server.CommitService

	chatID := update.FromChat().ID

	commit, err := commitService.FindByChatIDAndDate(chatID, state.data.SelectedDate)
	if err != nil {
		log.Error(
			"error when searching for a commit",
			"chatID", chatID, "state", CommitModify, "step", "Init", "err", err)

		gomer.SendMessage(chatID, i18n.Localize("oops"))

		state.stateMachine.
			Set(Idle, chatID, &StateContext{}).
			Handle(update)

		return
	}

	if commit == nil {
		msg := fmt.Sprintf(i18n.Localize("commitNotFound"), state.data.SelectedDate)

		gomer.SendMessage(chatID, msg)

		state.stateMachine.
			Set(Idle, chatID, &StateContext{}).
			Handle(update)

		return
	}

	state.data.Commit = commit

	msg := fmt.Sprintf(i18n.Localize("commitModifyPromt"), commit.Payload, commit.Date)

	gomer.SendMessage(chatID, msg)
}

func (state *CommitModifyState) Handle(update *tgbotapi.Update) {
	gomer := state.server.Gomer
	commitService := state.server.CommitService

	chatID := update.FromChat().ID

	// validate input
	if update.Message == nil {
		log.Error(
			"message is nil",
			"chatID", chatID, "state", Commit)

		state.Init(update)
		return
	}

	// extract message
	payload := update.Message.Text

	// update commit
	commit, err := commitService.UpdateCommit(state.data.Commit.ID, payload)
	if err != nil {
		log.Error(
			"error updating commit",
			"chatID", chatID, "err", err)

		gomer.SendMessage(chatID, i18n.Localize("commitModifyError"))

		state.stateMachine.
			Set(Idle, chatID, &StateContext{}).
			Init(update)

		return
	}

	// notify about success
	msg := fmt.Sprintf(i18n.Localize("commitModifySuccess"), commit.Date, commit.Payload)

	gomer.SendMessage(chatID, msg)

	// set starting state
	state.stateMachine.
		Set(Idle, chatID, &StateContext{}).
		Init(update)
}
