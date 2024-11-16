package state

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/epistax1s/gomer/internal/i18n"
	"github.com/epistax1s/gomer/internal/log"
	"github.com/epistax1s/gomer/internal/server"
)

const (
	cmdStart        = "start"
	cmdHelp         = "help"
	cmdTrack        = "track"
	cmdUntrack      = "untrack"
	cmdCommit       = "commit"
	cmdCommitModify = "modify"
)

type IdleCmdHandler func(server *server.Server, update *tgbotapi.Update)

var IdleCmdHandlers = map[string]IdleCmdHandler{
	cmdStart:        helpHandler,
	cmdHelp:         helpHandler,
	cmdTrack:        trackHandler,
	cmdUntrack:      untrackHandler,
	cmdCommit:       commitHandler,
	cmdCommitModify: commitModifyHandler,
}

func helpHandler(server *server.Server, update *tgbotapi.Update) {
	chatID := update.FromChat().ID

	userExists, err := server.UserService.UserExists(chatID)
	if err != nil {
		log.Error(
			"error checking the user's existence",
			"state", Idle, "cmd", cmdHelp, "err", err)

		return
	}

	server.Gomer.SendMessage(chatID, buildHelpResponse(userExists))
}

func buildHelpResponse(isTracked bool) string {
	if isTracked {
		return fmt.Sprintf(""+
			"/help		- %s\n"+
			"/commit 	- %s\n"+
			"/modify	- %s\n"+
			"/untrack	- %s",
			i18n.Localize("helpDescription"), i18n.Localize("commitDescription"),
			i18n.Localize("modifyDescription"), i18n.Localize("untrackDescription"))
	} else {
		return fmt.Sprintf(""+
			"/help	- %s\n"+
			"/track	- %s",
			i18n.Localize("helpDescription"), i18n.Localize("trackDescription"))
	}
}

func trackHandler(server *server.Server, update *tgbotapi.Update) {
	chatID := update.FromChat().ID

	userExists, err := server.UserService.UserExists(chatID)
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

	StateMachine.
		Set(TrackDepartment, chatID, &StateContext{}).
		Init(server, update)
}

func untrackHandler(server *server.Server, update *tgbotapi.Update) {
	chatID := update.FromChat().ID

	userExists, err := server.UserService.UserExists(chatID)
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

	if err := server.UserService.UntrackUser(chatID); err != nil {
		log.Error(
			"error when trying to stop tracking the user",
			"state", Idle, "cmd", cmdUntrack, "chatID", chatID, "err", err)

		server.Gomer.SendMessage(chatID, i18n.Localize("oops"))
		return
	}

	if err := server.UserService.UntrackUser(chatID); err == nil {
		server.Gomer.SendMessage(chatID, i18n.Localize("trackStopped"))
	} else {
		server.Gomer.SendMessage(chatID, i18n.Localize("oops"))
	}
}

func commitHandler(server *server.Server, update *tgbotapi.Update) {
	chatID := update.FromChat().ID

	userExists, err := server.UserService.UserExists(chatID)
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

	StateMachine.
		Set(Date, chatID, &StateContext{NextState: Commit}).
		Init(server, update)
}

func commitModifyHandler(server *server.Server, update *tgbotapi.Update) {
	chatID := update.FromChat().ID

	userExists, err := server.UserService.UserExists(chatID)
	if err != nil {
		log.Error(
			"error checking the user's existence",
			"state", Idle, "cmd", cmdCommitModify, "err", err)
	}

	if !userExists {
		log.Info(
			"the user is not a tracked",
			"chatID", chatID)

		return
	}

	StateMachine.
		Set(Date, chatID, &StateContext{NextState: CommitModify}).
		Init(server, update)
}
