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

	if update.Message != nil && update.Message.Command() == "start" {
		log.Info("i.Register(update)")
		i.Register(update)
	} else {
		log.Info("i.IsRegistered(update)")
		i.IsRegistered(update)
	}
}

func (i *SecurityInterceptor) Register(update *tgbotapi.Update) {
	gomer := i.Server.Gomer
	authUserService := i.Server.AuthUserService
	authKeyService := i.Server.AuthKeyService

	chatID := update.FromChat().ID
	username := update.FromChat().UserName

	key := update.Message.CommandArguments()

	log.Info("attempt to register a user.", "chatID", chatID, "username", username, "key", key)

	if authUserService.IsRegistered(chatID) {
		log.Info("the user is already registered", "chatID", chatID, "username", username, "key", key)
		i.Next(update)
	} else {
		if !authKeyService.IsValidKey(key) {
			log.Error("invalid key", "chatID", chatID, "username", username, "key", key)
			gomer.SendMessage(chatID, i18n.Localize("accessDenied"))
			return
		}

		if err := authUserService.Register(chatID, username); err == nil {
			log.Info("the user has been successfully registered", "chatID", chatID, "username", username, "key", key)
			i.Next(update)
		} else {
			log.Error("error during user registration", "chatID", chatID, "username", username, "key", key, "err", err)
			gomer.SendMessage(chatID, i18n.Localize("accessDenied"))
		}
	}
}

func (i *SecurityInterceptor) IsRegistered(update *tgbotapi.Update) {
	gomer := i.Server.Gomer
	authUserService := i.Server.AuthUserService

	chatID := update.FromChat().ID
	username := update.FromChat().UserName

	log.Info("verifying user access", "chatID", chatID, "username", username)

	if authUserService.IsRegistered(chatID) {
		log.Info("the user is authenticated, access is allowed", "chatID", chatID, "username", username)
		i.Next(update)
	} else {
		log.Info("user authentication error, access is denied", "chatID", chatID, "username", username)
		gomer.SendMessage(chatID, i18n.Localize("accessDenied"))
	}
}
