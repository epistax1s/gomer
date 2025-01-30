package manage_users

import (
	callback "github.com/epistax1s/gomer/internal/callback"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const USERS_PAGE_SIZE = 10

func (state *ManageUsersState) Init(update *tgbotapi.Update) {
}

func (state *ManageUsersState) Handle(update *tgbotapi.Update) {
}

func (state *ManageUsersState) prevHandler(update *tgbotapi.Update, c callback.Callback) {
}

func (state *ManageUsersState) nextHandler(update *tgbotapi.Update, c callback.Callback) {
}

func (state *ManageUsersState) selectHandler(update *tgbotapi.Update, c callback.Callback) {
}
