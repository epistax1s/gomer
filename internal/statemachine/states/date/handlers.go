package date

import (
	"time"

	"github.com/epistax1s/gomer/internal/calendar"
	"github.com/epistax1s/gomer/internal/i18n"
	"github.com/epistax1s/gomer/internal/log"

	. "github.com/epistax1s/gomer/internal/statemachine/core"

	callback "github.com/epistax1s/gomer/internal/callback"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (state *DateState) Init(update *tgbotapi.Update) {
	gomer := state.server.Gomer

	chatID := update.FromChat().ID
	now := time.Now()

	calendarMarkup := calendar.GenerateCalendar(now.Year(), now.Month())

	gomer.SendMessageWithKeyboard(
		chatID, i18n.Localize("chooseDatePromt"), calendarMarkup)
}

func (state *DateState) Handle(update *tgbotapi.Update) {
	gomer := state.server.Gomer

	chatID := update.FromChat().ID
	query := update.CallbackQuery

	if query != nil {
		queryData := query.Data
		callback, err := callback.Decode(queryData)
		if err != nil {
			log.Error(err.Error())
		}
		state.handlers[callback.GetType()](update, callback)
	} else {
		gomer.SendMessage(chatID, i18n.Localize("oops"))
	}
}

func (state *DateState) prevHandler(update *tgbotapi.Update, c callback.Callback) {
	gomer := state.server.Gomer

	chatID := update.CallbackQuery.From.ID
	messageID := update.CallbackQuery.Message.MessageID

	prevCallback := c.(*callback.CalendarPrevCallback)
	calendarMarkup, _, _ := calendar.HandlePrevButton(prevCallback.Year, prevCallback.Month)

	gomer.EditMessageWithKeyboard(chatID, messageID, i18n.Localize("chooseDatePromt"), calendarMarkup)
}

func (state *DateState) nextHandler(update *tgbotapi.Update, c callback.Callback) {
	gomer := state.server.Gomer

	chatID := update.CallbackQuery.From.ID
	messageID := update.CallbackQuery.Message.MessageID

	nextCallback := c.(*callback.CalendarNextCallback)
	calendarMarkup, _, _ := calendar.HandleNextButton(nextCallback.Year, nextCallback.Month)

	gomer.EditMessageWithKeyboard(chatID, messageID, i18n.Localize("chooseDatePromt"), calendarMarkup)
}

func (state *DateState) dateHandler(update *tgbotapi.Update, c callback.Callback) {
	gomer := state.server.Gomer

	chatID := update.CallbackQuery.From.ID

	dateCallback := c.(*callback.CalendarDateCallback)
	date := dateCallback.Date

	commitDate, err := calendar.HandleButtonData(date)
	if err != nil {
		log.Error(
			"error parsing callback data",
			"chatID", chatID, "state", Date, "step", "Handle", "callback", c, "err", err)

		gomer.SendMessage(chatID, i18n.Localize("oops"))

		state.stateMachine.
			Set(Date, chatID, &StateContext{}).
			Init(update)

		return
	}

	gomer.SendCallbackResponse(update.CallbackQuery, i18n.Localize("chooseDateSuccess"))

	nextState := state.data.NextState
	log.Info(
		"the user has selected a date",
		"chatID", chatID, "state", Date, "date", commitDate, "nextState", nextState)

	state.stateMachine.
		Set(nextState, chatID, &StateContext{SelectedDate: commitDate}).
		Init(update)
}
