package state

import (
	"time"

	"github.com/epistax1s/gomer/internal/calendar"
	"github.com/epistax1s/gomer/internal/i18n"
	"github.com/epistax1s/gomer/internal/log"
	"github.com/epistax1s/gomer/internal/server"

	callback "github.com/epistax1s/gomer/internal/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type DateState struct {
	data *StateContext
}

func NewDateState(data *StateContext) State {
	return &DateState{
		data: data,
	}
}

type dateCallbackHandler func(*server.Server, *tgbotapi.Update, callback.Callback, *StateContext)

var dateCallbacHandlers = map[string]dateCallbackHandler{
	callback.CALENDAR_PREV: calendarPrevHandler,
	callback.CALENDAR_NEXT: calendarNextHandler,
	callback.CALENDAR_DATE: calendarDateHandler,
}

func (state *DateState) Init(server *server.Server, update *tgbotapi.Update) {
	chatID := update.FromChat().ID
	now := time.Now()

	calendarMarkup := calendar.GenerateCalendar(now.Year(), now.Month())

	server.Gomer.SendMessageWithKeyboard(
		chatID, i18n.Localize("chooseDatePromt"), calendarMarkup)
}

func (state *DateState) Handle(server *server.Server, update *tgbotapi.Update) {
	chatID := update.FromChat().ID
	query := update.CallbackQuery

	if query != nil {
		queryData := query.Data
		callback, err := callback.Decode(queryData)
		if err != nil {
			log.Error(err.Error())
		}
		dateCallbacHandlers[callback.GetType()](server, update, callback, state.data)
	} else {
		server.Gomer.SendMessage(chatID, i18n.Localize("oops"))
	}
}

func calendarPrevHandler(server *server.Server, update *tgbotapi.Update, c callback.Callback, sc *StateContext) {
	chatID := update.CallbackQuery.From.ID
	messageID := update.CallbackQuery.Message.MessageID

	prevCallback := c.(*callback.CalendarPrevCallback)
	calendarMarkup, _, _ := calendar.HandlePrevButton(prevCallback.Year, prevCallback.Month)

	server.Gomer.EditMessageWithKeyboard(
		chatID, messageID, i18n.Localize("chooseDatePromt"), calendarMarkup)
}

func calendarNextHandler(server *server.Server, update *tgbotapi.Update, c callback.Callback, sc *StateContext) {
	chatID := update.CallbackQuery.From.ID
	messageID := update.CallbackQuery.Message.MessageID

	nextCallback := c.(*callback.CalendarNextCallback)
	calendarMarkup, _, _ := calendar.HandleNextButton(nextCallback.Year, nextCallback.Month)

	server.Gomer.EditMessageWithKeyboard(
		chatID, messageID, i18n.Localize("chooseDatePromt"), calendarMarkup)
}

func calendarDateHandler(server *server.Server, update *tgbotapi.Update, c callback.Callback, sc *StateContext) {
	chatID := update.CallbackQuery.From.ID

	dateCallback := c.(*callback.CalendarDateCallback)
	date := dateCallback.Date

	commitDate, err := calendar.HandleButtonData(date)
	if err != nil {
		log.Error(
			"error parsing callback data",
			"chatID", chatID, "state", Date, "step", "Handle", "callback", c, "err", err)

		server.Gomer.SendMessage(chatID, i18n.Localize("oops"))

		StateMachine.
			Set(Date, chatID, &StateContext{}).
			Init(server, update)

		return
	}

	server.Gomer.SendCallbackResponse(update.CallbackQuery, i18n.Localize("chooseDateSuccess"))

	log.Info(
		"the user has selected a date",
		"chatID", chatID, "state", Date, "date", commitDate, "nextState", sc.NextState)

	StateMachine.
		Set(sc.NextState, chatID, &StateContext{CommitDate: commitDate}).
		Init(server, update)
}
