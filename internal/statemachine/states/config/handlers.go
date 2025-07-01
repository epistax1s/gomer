package config

import (
	"fmt"

	callback "github.com/epistax1s/gomer/internal/callback"
	"github.com/epistax1s/gomer/internal/i18n"
	"github.com/epistax1s/gomer/internal/log"
	"github.com/epistax1s/gomer/internal/model"

	. "github.com/epistax1s/gomer/internal/statemachine/core"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (state *ConfigState) Init(update *tgbotapi.Update) {
	gomer := state.server.Gomer
	chatID := update.FromChat().ID

	msg := fmt.Sprintf(""+
		"/manual - %s\n"+
		"/redmine - %s\n"+
		"/redmine_ext - %s\n"+
		"/cancel - %s\n",
		i18n.Localize("commitSrcManualDescription"),
		i18n.Localize("commitSrcRedmineDescription"),
		i18n.Localize("commitSrcRedmineExtDescription"),
		i18n.Localize("commitSrcCancelDescription"))

	gomer.SendMessage(chatID, msg)
}

func (state *ConfigState) Handle(update *tgbotapi.Update) {
	gomer := state.server.Gomer

	chatID := update.FromChat().ID
	cmd := update.Message.Command()

	cmdHandler, exits := state.handlers[cmd]
	if exits {
		cmdHandler(update, nil)
	} else {
		gomer.SendMessage(chatID, i18n.Localize("unsupportedCommand"))
		state.Init(update)
	}
}

func (state *ConfigState) manualHandler(update *tgbotapi.Update, callback callback.Callback) {
	state.changeCommitSrc(update, model.UserCommitSrcManual)
}

func (state *ConfigState) redmineHandler(update *tgbotapi.Update, callback callback.Callback) {
	state.changeCommitSrc(update, model.UserCommitSrcRedmine)
}

func (state *ConfigState) redmineExtHandler(update *tgbotapi.Update, callback callback.Callback) {
	state.changeCommitSrc(update, model.UserCommitSrcRedmineExt)
}

func (state *ConfigState) changeCommitSrc(update *tgbotapi.Update, commitSrc string) {
	gomer := state.server.Gomer
	userService := state.server.UserService

	chatID := update.FromChat().ID

	if err := userService.SetCommitSrc(chatID, commitSrc); err != nil {
		log.Error(
			"error trying to change commitSrc for user",
			"chatID", chatID, "state", Config, "commitSrc", commitSrc, "err", err)

		gomer.SendMessage(chatID, i18n.Localize("oops"))

		state.stateMachine.
			Set(Idle, chatID, &StateContext{}).
			Handle(update)

		return
	}

	log.Info(
		"commitSrc changed successfully",
		"chatID", chatID, "commitSrc", commitSrc)

	msg := fmt.Sprintf(i18n.Localize("commitSrcChangedSuccess"), commitSrc)

	gomer.SendMessage(chatID, msg)

	state.stateMachine.
		Set(Idle, chatID, &StateContext{}).
		Init(update)
}
