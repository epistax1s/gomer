package state

import (
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/epistax1s/gomer/internal/calendar"
	"github.com/epistax1s/gomer/internal/i18n"
	"github.com/epistax1s/gomer/internal/log"
	"github.com/epistax1s/gomer/internal/server"
)

type DateState struct {
	data *StateContext
}

func NewDateState(data *StateContext) State {
	return &DateState{
		data: data,
	}
}

func (state *DateState) Init(server *server.Server, update *tgbotapi.Update) {
	chatID := update.FromChat().ID
	now := time.Now()

	calendarMarkup := calendar.GenerateCalendar(now.Year(), now.Month())

	server.Gomer.SendMessageWithMarkup(
		chatID, i18n.Localize("chooseDatePromt"), calendarMarkup)
}

func (state *DateState) Handle(server *server.Server, update *tgbotapi.Update) {
	chatID := update.FromChat().ID

	if update.CallbackQuery == nil {
		log.Warn(
			"callbackQuery is nil",
			"chatID", chatID, "state", Date, "step", "Handle")

		server.Gomer.SendMessage(chatID, i18n.Localize("chooseDatePromt"))
		return
	}

	// extracting the selected date
	data := update.CallbackQuery.Data
	commitDate, err := calendar.HandleButtonData(update.CallbackQuery.Data)
	if err != nil {
		log.Error(
			"error parsing callback data",
			"chatID", chatID, "state", Date, "step", "Handle", "callbackData", data, "err", err)

		server.Gomer.SendMessage(chatID, i18n.Localize("oops"))
		state.Init(server, update)

		return
	}

	callbackConfig := tgbotapi.CallbackConfig{
		CallbackQueryID: update.CallbackQuery.ID,
		Text:            i18n.Localize("chooseDateSuccess"),
		ShowAlert:       false,
	}

	if _, err := server.Gomer.BotAPI.Request(callbackConfig); err != nil {
		log.Error(
			"error confirming callback processing",
			"chatID", chatID, "state", Date, "err", err)

		return
	}

	log.Info(
		"the user has selected a date",
		"chatID", chatID, "state", Date, "date", commitDate, "nextState", state.data.NextState)

	StateMachine.
		Set(state.data.NextState, chatID, &StateContext{CommitDate: commitDate}).
		Init(server, update)
}
