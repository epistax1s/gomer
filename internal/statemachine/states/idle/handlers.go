package idle

import (
	"fmt"
	"strings"

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
	securityService := state.server.SecurityService

	chatID := update.FromChat().ID

	isActive := securityService.IsActive(chatID)
	isAdmin := securityService.IsAdmin(chatID)

	var help string

	if isActive {
		help = fmt.Sprintf(""+
			"/help - %s\n"+
			"/commit - %s\n"+
			"/modify - %s\n"+
			"/config - %s\n"+
			"/untrack - %s\n"+
			"/invites - %s\n"+
			"/invite_new - %s\n"+
			"/cancel - %s",
			i18n.Localize("helpDescription"),
			i18n.Localize("commitDescription"),
			i18n.Localize("modifyDescription"),
			i18n.Localize("configDescription"),
			i18n.Localize("untrackDescription"),
			i18n.Localize("myInvitesDescription"),
			i18n.Localize("generateInviteDescription"),
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
	gomer := state.server.Gomer
	securityService := state.server.SecurityService

	chatID := update.FromChat().ID

	if isRegistered := securityService.IsRegistered(chatID); !isRegistered {
		log.Error(
			"this user is not registered",
			"state", Idle, "cmd", cmdTrack, "chatID", chatID)

		gomer.SendMessage(chatID, i18n.Localize("oops"))
		return
	}

	if isActive := securityService.IsActive(chatID); isActive {
		log.Error(
			"the user is already being tracked",
			"chatID", chatID)

		gomer.SendMessage(chatID, i18n.Localize("oops"))
		return
	}

	state.stateMachine.
		Set(TrackDepartment, chatID, &StateContext{}).
		Init(update)
}

func (state *IdleState) untrackHandler(update *tgbotapi.Update, callback callback.Callback) {
	gomer := state.server.Gomer
	userService := state.server.UserService
	securityService := state.server.SecurityService

	chatID := update.FromChat().ID

	if isActive := securityService.IsActive(chatID); !isActive {
		log.Info(
			"the user is not active",
			"chatID", chatID)

		gomer.SendMessage(chatID, i18n.Localize("oops"))
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
	gomer := state.server.Gomer
	securityService := state.server.SecurityService

	chatID := update.FromChat().ID

	if isActive := securityService.IsActive(chatID); !isActive {
		log.Info(
			"the user is not active",
			"chatID", chatID)
		gomer.SendMessage(chatID, i18n.Localize("oops"))
		return
	}

	state.stateMachine.
		Set(Date, chatID, &StateContext{NextState: Commit}).
		Init(update)
}

func (state *IdleState) commitModifyHandler(update *tgbotapi.Update, callback callback.Callback) {
	gomer := state.server.Gomer
	securityService := state.server.SecurityService

	chatID := update.FromChat().ID

	if isActive := securityService.IsActive(chatID); !isActive {
		log.Info(
			"the user is not active",
			"chatID", chatID)

		gomer.SendMessage(chatID, i18n.Localize("oops"))
		return
	}

	state.stateMachine.
		Set(Date, chatID, &StateContext{NextState: CommitModify}).
		Init(update)
}

func (state *IdleState) configHandler(update *tgbotapi.Update, callback callback.Callback) {
	gomer := state.server.Gomer
	securityService := state.server.SecurityService

	chatID := update.FromChat().ID

	if isActive := securityService.IsActive(chatID); !isActive {
		log.Info(
			"the user is not active",
			"chatID", chatID)

		gomer.SendMessage(chatID, i18n.Localize("oops"))
		return
	}

	state.stateMachine.
		Set(Config, chatID, &StateContext{}).
		Init(update)
}

func (state *IdleState) forcePublishHandler(update *tgbotapi.Update, callback callback.Callback) {
	gomer := state.server.Gomer
	securityService := state.server.SecurityService

	chatID := update.FromChat().ID

	if isAdmin := securityService.IsAdmin(chatID); !isAdmin {
		log.Info(
			"the user is not a admin",
			"chatID", chatID)

		gomer.SendMessage(chatID, i18n.Localize("oops"))
		return
	}

	state.stateMachine.
		Set(Date, chatID, &StateContext{NextState: ForcePublish}).
		Init(update)
}

func (state *IdleState) manageUsersHandler(update *tgbotapi.Update, callback callback.Callback) {
	gomer := state.server.Gomer
	securityService := state.server.SecurityService

	chatID := update.FromChat().ID

	if isAdmin := securityService.IsAdmin(chatID); !isAdmin {
		log.Info(
			"the user is not a admin",
			"chatID", chatID)

		gomer.SendMessage(chatID, i18n.Localize("oops"))
		return
	}

	state.stateMachine.
		Set(ManageUsers, chatID, &StateContext{NextState: CommitModify}).
		Init(update)
}

func (state *IdleState) manageGroupsHandler(update *tgbotapi.Update, callback callback.Callback) {
	gomer := state.server.Gomer
	securityService := state.server.SecurityService

	chatID := update.FromChat().ID

	if isAdmin := securityService.IsAdmin(chatID); !isAdmin {
		log.Info(
			"the user is not a admin",
			"chatID", chatID)

		gomer.SendMessage(chatID, i18n.Localize("oops"))
		return
	}

	state.stateMachine.
		Set(ManageGrops, chatID, &StateContext{NextState: CommitModify}).
		Init(update)
}

func (state *IdleState) generateInviteHandler(update *tgbotapi.Update, callback callback.Callback) {
	securityService := state.server.SecurityService
	gomer := state.server.Gomer

	chatID := update.FromChat().ID

	if isActive := securityService.IsActive(chatID); !isActive {
		log.Info(
			"the user is not active",
			"chatID", chatID)

		gomer.SendMessage(chatID, i18n.Localize("oops"))
		return
	}

	inviteLink, err := securityService.GenerateInviteLink(chatID)
	if err != nil {
		log.Error("failed to generate invite link", "error", err)
		gomer.SendMessage(chatID, i18n.Localize("inviteGenerationFailed"))
		return
	}

	gomer.SendMessage(chatID, fmt.Sprintf(i18n.Localize("inviteLinkGenerated"), inviteLink))
}

func (state *IdleState) myInvitesHandler(update *tgbotapi.Update, callback callback.Callback) {
	securityService := state.server.SecurityService
	gomer := state.server.Gomer

	chatID := update.FromChat().ID

	if isActive := securityService.IsActive(chatID); !isActive {
		log.Info(
			"the user is not active",
			"chatID", chatID)

		gomer.SendMessage(chatID, i18n.Localize("oops"))
		return
	}

	invites, err := securityService.GetUserInvites(chatID)
	if err != nil {
		log.Error("failed to get user invites", "error", err)
		gomer.SendMessage(chatID, i18n.Localize("oops"))
		return
	}

	if len(invites) == 0 {
		gomer.SendMessage(chatID, i18n.Localize("noInvites"))
		return
	}

	var message strings.Builder
	message.WriteString(i18n.Localize("yourInvites_html"))

	botUsername := state.server.Config.Bot.Username
	for _, invite := range invites {
		link := fmt.Sprintf("https://t.me/%s?start=%s", botUsername, invite.Code)
		
		message.WriteString("\n\n")
		message.WriteString(fmt.Sprintf(i18n.Localize("inviteLink_html"), link))
		message.WriteString("\n")
		message.WriteString(fmt.Sprintf(i18n.Localize("inviteCreatedAt_html"), invite.CreatedAt))
		message.WriteString("\n")

		if invite.Used {
			message.WriteString(fmt.Sprintf(i18n.Localize("inviteUsedBy_html"), invite.UsedBy.Name))
		} else {
			message.WriteString(i18n.Localize("inviteNotUsed_html"))
		}
	}

	gomer.SendMessageHtml(chatID, message.String())
}
