package idle

import (
	"fmt"

	"github.com/epistax1s/gomer/internal/callback"
	. "github.com/epistax1s/gomer/internal/statemachine/core"

	"github.com/epistax1s/gomer/internal/i18n"
	"github.com/epistax1s/gomer/internal/log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (state *IdleState) Init(update *tgbotapi.Update) {
	helpHandler := state.handlers[cmdHelp]
	helpHandler(update, nil)
}

func (state *IdleState) Handle(update *tgbotapi.Update) {
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

func (state *IdleState) helpHandler(update *tgbotapi.Update, callback callback.Callback) {
	gomer := state.server.Gomer

	userService := state.server.UserService
	securityService := state.server.SecurityService

	chatID := update.FromChat().ID

	userExists, err := userService.UserExists(chatID)
	if err != nil {
		log.Error(
			" error when checking user existence",
			"state", Idle, "cmd", cmdHelp, "err", err)

		return
	}

	isAdmin, err := securityService.IsAdmin(chatID)
	if err != nil {
		log.Error(
			"user role validation error",
			"state", Idle, "cmd", cmdHelp, "err", err)

		return
	}

	var help string

	if userExists {
		help = fmt.Sprintf(""+
			"/help		- %s\n"+
			"/commit 	- %s\n"+
			"/modify	- %s\n"+
			"/untrack	- %s\n"+
			"/cancel	- %s",
			i18n.Localize("helpDescription"), i18n.Localize("commitDescription"),
			i18n.Localize("modifyDescription"), i18n.Localize("untrackDescription"),
			i18n.Localize("cancelDescription"),
		)
		if isAdmin {
			help = fmt.Sprintf(""+
				"%s\n"+
				"/publish	- %s",
				help, i18n.Localize("publishDescription"),
			)
		}

	} else {
		help = fmt.Sprintf(""+
			"/help		- %s\n"+
			"/track		- %s\n"+
			"/cancel	- %s",
			i18n.Localize("helpDescription"), i18n.Localize("trackDescription"),
			i18n.Localize("cancelDescription"),
		)
	}

	gomer.SendMessage(chatID, help)
}

func (state *IdleState) trackHandler(update *tgbotapi.Update, callback callback.Callback) {
	userService := state.server.UserService

	chatID := update.FromChat().ID

	userExists, err := userService.UserExists(chatID)
	if err != nil {
		log.Error(
			"error checking the user's existence",
			"state", Idle, "cmd", cmdTrack, "err", err)

		return
	}

	if userExists {
		log.Info(
			"the user is already being tracked",
			"chatID", chatID)

		return
	}

	state.stateMachine.
		Set(TrackDepartment, chatID, &StateContext{}).
		Init(update)
}

func (state *IdleState) untrackHandler(update *tgbotapi.Update, callback callback.Callback) {
	gomer := state.server.Gomer
	userService := state.server.UserService

	chatID := update.FromChat().ID

	userExists, err := userService.UserExists(chatID)
	if err != nil {
		log.Error(
			"error checking the user's existence",
			"state", Idle, "cmd", cmdUntrack, "err", err)

		return
	}

	if !userExists {
		log.Info(
			"the user is not a tracked",
			"chatID", chatID)

		return
	}

	if err := userService.UntrackUser(chatID); err != nil {
		log.Error(
			"error when trying to stop tracking the user",
			"state", Idle, "cmd", cmdUntrack, "chatID", chatID, "err", err)

		gomer.SendMessage(chatID, i18n.Localize("oops"))
		return
	}

	if err := userService.UntrackUser(chatID); err == nil {
		gomer.SendMessage(chatID, i18n.Localize("trackStopped"))
	} else {
		gomer.SendMessage(chatID, i18n.Localize("oops"))
	}
}

func (state *IdleState) commitHandler(update *tgbotapi.Update, callback callback.Callback) {
	userService := state.server.UserService

	chatID := update.FromChat().ID

	userExists, err := userService.UserExists(chatID)
	if err != nil {
		log.Error(
			"error checking the user's existence",
			"state", Idle, "cmd", cmdCommit, "err", err)
		return
	}

	if !userExists {
		log.Info(
			"the user is not a tracked",
			"chatID", chatID)

		return
	}

	state.stateMachine.
		Set(Date, chatID, &StateContext{NextState: Commit}).
		Init(update)
}

func (state *IdleState) commitModifyHandler(update *tgbotapi.Update, callback callback.Callback) {
	userService := state.server.UserService

	chatID := update.FromChat().ID

	userExists, err := userService.UserExists(chatID)
	if err != nil {
		log.Error(
			"error checking the user's existence",
			"state", Idle, "cmd", cmdCommitModify, "err", err)
		return
	}

	if !userExists {
		log.Info(
			"the user is not a tracked",
			"chatID", chatID)

		return
	}

	state.stateMachine.
		Set(Date, chatID, &StateContext{NextState: CommitModify}).
		Init(update)
}

func (state *IdleState) forcePublishHandler(update *tgbotapi.Update, callback callback.Callback) {
	securityService := state.server.SecurityService

	chatID := update.FromChat().ID

	isAdmin, err := securityService.IsAdmin(chatID)
	if err != nil {
		log.Error(
			"Force publish handler error",
			"err", err.Error())
		return
	}

	if !isAdmin {
		log.Warn(
			"Force publish is available only to the administrator",
			"chatID", chatID)
		return
	}

	state.stateMachine.
		Set(Date, chatID, &StateContext{NextState: ForcePublish}).
		Init(update)
}

func (state *IdleState) manageUsersHandler(update *tgbotapi.Update, callback callback.Callback) {
	securityService := state.server.SecurityService

	chatID := update.FromChat().ID

	isAdmin, err := securityService.IsAdmin(chatID)
	if err != nil {
		log.Error(
			"User management handler error",
			"err", err.Error())
		return
	}

	if !isAdmin {
		log.Warn(
			"User management is available only to the administrator",
			"chatID", chatID)
		return
	}

	state.stateMachine.
		Set(ManageUsers, chatID, &StateContext{NextState: CommitModify}).
		Init(update)
}

func (state *IdleState) manageGroupsHandler(update *tgbotapi.Update, callback callback.Callback) {
	securityService := state.server.SecurityService

	chatID := update.FromChat().ID

	isAdmin, err := securityService.IsAdmin(chatID)
	if err != nil {
		log.Error(
			"Group management handler error",
			"err", err.Error())
		return
	}

	if !isAdmin {
		log.Warn(
			"Group management is available only to the administrator",
			"chatID", chatID)
		return
	}

	state.stateMachine.
		Set(ManageGrops, chatID, &StateContext{NextState: CommitModify}).
		Init(update)
}
