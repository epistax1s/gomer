package interceptor

import (
	"github.com/epistax1s/gomer/internal/i18n"
	"github.com/epistax1s/gomer/internal/log"
	"github.com/epistax1s/gomer/internal/server"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type SecurityInterceptor struct {
	Server *server.Server
	BaseInterceptor
}

func (i *SecurityInterceptor) Handle(update *tgbotapi.Update) {
	if !update.FromChat().IsPrivate() {
		i.Next(update)
		return
	}

	gomer := i.Server.Gomer

	chat :=update.FromChat()

	chatID := chat.ID
	username := chat.UserName
	name := chat.FirstName + chat.LastName

	securityService := i.Server.SecurityService

	if update.Message != nil && update.Message.Command() == "start" {
		// Registration flow with invite code
		inviteCode := update.Message.CommandArguments()
		if inviteCode == "" {
			gomer.SendMessage(chatID, i18n.Localize("inviteCodeRequired"))
			return
		}

		err := securityService.RegisterUser(chatID, username, name, inviteCode)
		if err != nil {
			log.Error("failed to register user", "chatID", chatID, "username", username, "error", err)
			gomer.SendMessage(chatID, i18n.Localize("registrationFailed"))
			return
		}

		log.Info("user successfully registered", "chatID", chatID, "username", username)
		gomer.SendMessage(chatID, i18n.Localize("registrationSuccess"))
		i.Next(update)
		return
	}

	// Authentication flow for existing users
	isAuthenticated, err := securityService.AuthenticateUser(chatID)
	if err != nil {
		log.Error("authentication error", "chatID", chatID, "username", username, "error", err)
		gomer.SendMessage(chatID, i18n.Localize("authError"))
		return
	}

	if !isAuthenticated {
		log.Info("unauthorized access attempt", "chatID", chatID, "username", username)
		gomer.SendMessage(chatID, i18n.Localize("accessDenied"))
		return
	}

	i.Next(update)
}
