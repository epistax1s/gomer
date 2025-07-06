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
	userService := state.server.UserService
	chatID := update.FromChat().ID

	user, _ := userService.FindUserByChatID(chatID)
	userCommitSrc := fmt.Sprintf(i18n.Localize("profileCommitSrc_html"), user.CommitSrc)
	userRedmineID := fmt.Sprintf(i18n.Localize("profileRedmineID_html"), user.RedmineID)

	msg := fmt.Sprintf(""+
		"%s\n\n"+ // Profile title
		"%s\n"+
		"%s\n\n"+
		"%s\n\n"+ // Commit source title
		"/manual - %s\n\n"+
		"/redmine - %s\n\n"+
		"/redmine_ext - %s\n\n"+
		"/examples - %s\n\n"+
		"%s\n\n"+ // Redmine id title
		"/redmine_id - %s\n\n"+
		"%s\n\n"+ // Other title
		"/cancel - %s\n",
		i18n.Localize("profileTitle_html"),
		userCommitSrc,
		userRedmineID,
		i18n.Localize("commitSrcTitle_html"),
		i18n.Localize("commitSrcManualDescription"),
		i18n.Localize("commitSrcRedmineDescription"),
		i18n.Localize("commitSrcRedmineExtDescription"),
		i18n.Localize("commitSrcExamplesDescription"),
		i18n.Localize("redmineIdTitle_html"),
		i18n.Localize("redmineIdDescription"),
		i18n.Localize("confOtherTitle_html"),
		i18n.Localize("confCancelDescription"),
	)

	gomer.SendMessageHtml(chatID, msg)
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

func (state *ConfigState) examplesHandler(update *tgbotapi.Update, callback callback.Callback) {
	state.server.Gomer.SendMessageHtml(
		update.FromChat().ID,
		i18n.Localize("commitSrcExamples_html"),
	)
}

func (state *ConfigState) changeCommitSrc(update *tgbotapi.Update, commitSrc string) {
	gomer := state.server.Gomer
	userService := state.server.UserService

	chatID := update.FromChat().ID

	redmineSrc := commitSrc == model.UserCommitSrcRedmine || commitSrc == model.UserCommitSrcRedmineExt
	user, _ := userService.FindUserByChatID(chatID)
	if redmineSrc && user.RedmineID == nil {
		gomer.SendMessage(chatID, i18n.Localize("redmineIdRequired"))
		state.Init(update)
		return
	}

	user.CommitSrc = commitSrc
	userService.Save(user)

	log.Info(
		"commitSrc changed successfully",
		"chatID", chatID, "commitSrc", commitSrc)

	msg := fmt.Sprintf(i18n.Localize("commitSrcChangedSuccess"), commitSrc)

	gomer.SendMessage(chatID, msg)

	state.stateMachine.
		Set(Idle, chatID, &StateContext{}).
		Init(update)
}

func (state *ConfigState) redmineID(update *tgbotapi.Update, callback callback.Callback) {
	state.stateMachine.
		Set(RedmineID, update.FromChat().ID, &StateContext{}).
		Init(update)
}
