package commit

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/epistax1s/gomer/internal/i18n"
	"github.com/epistax1s/gomer/internal/log"

	. "github.com/epistax1s/gomer/internal/statemachine/core"
)

func (state *CommitState) Init(update *tgbotapi.Update) {
	gomer := state.server.Gomer
	commitService := state.server.CommitService

	chatID := update.FromChat().ID

	commit, err := commitService.FindCommitByUserIdAndDate(chatID, state.data.CommitDate)
	if err != nil {
		log.Error(
			"error when searching for a commit",
			"chatID", chatID, "state", commit, "step", "Init", "err", err)

		gomer.SendMessage(chatID, i18n.Localize("oops"))

		state.stateMachine.
			Set(Idle, chatID, &StateContext{}).
			Handle(update)

		return
	}

	if commit != nil {
		log.Info(
			"attempt to create a second commit on the same date",
			"chatID", chatID, "state", commit, "step", "Init", "date", state.data.CommitDate)

		msg := fmt.Sprintf(i18n.Localize("commitAlreadyCreated"), state.data.CommitDate, commit.Payload)

		gomer.SendMessage(chatID, msg)

		state.stateMachine.
			Set(Idle, chatID, &StateContext{}).
			Init(update)

		return
	}

	msg := fmt.Sprintf(i18n.Localize("commitCreatePromt"), state.data.CommitDate)

	gomer.SendMessage(chatID, msg)
}

func (state *CommitState) Handle(update *tgbotapi.Update) {
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

	payload := update.Message.Text

	// create commit
	err := commitService.CreateCommit(chatID, payload, state.data.CommitDate)
	if err != nil {
		gomer.SendMessage(chatID, i18n.Localize("commitCreateError"))
		return
	}

	// send success response
	msg := fmt.Sprintf(i18n.Localize("commitCreateSuccess"), state.data.CommitDate, payload)

	gomer.SendMessage(chatID, msg)

	// set base state
	state.stateMachine.
		Set(Idle, chatID, &StateContext{}).
		Init(update)
}
