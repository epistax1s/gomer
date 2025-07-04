package publish

import (
	"fmt"
	"time"

	"github.com/epistax1s/gomer/internal/callback"
	"github.com/epistax1s/gomer/internal/database"
	"github.com/epistax1s/gomer/internal/i18n"
	"github.com/epistax1s/gomer/internal/log"

	. "github.com/epistax1s/gomer/internal/statemachine/core"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (state *ForcePublishState) Init(update *tgbotapi.Update) {
	gomer := state.server.Gomer
	chatID := update.FromChat().ID

	confirmKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(i18n.Localize("action_confirm"), callback.NewActionConfirmCallback()),
			tgbotapi.NewInlineKeyboardButtonData(i18n.Localize("action_cancel"), callback.NewActionCancelCallback()),
		),
	)
	confirmMsg := fmt.Sprintf(i18n.Localize("forcePublishConfirm"), state.data.SelectedDate.String())
	gomer.SendMessageWithKeyboard(chatID, confirmMsg , &confirmKeyboard)
}

func (state *ForcePublishState) Handle(update *tgbotapi.Update) {
	gomer := state.server.Gomer

	chatID := update.FromChat().ID
	query := update.CallbackQuery

	if query != nil {
		// confirm receipt CallbackQuery
		gomer.SendCallbackResponse(query, i18n.Localize("ok"))

		callback, err := callback.Decode(query.Data)
		if err != nil {
			log.Error(err.Error())
			gomer.SendMessage(chatID, i18n.Localize("oops"))
		}
		state.handlers[callback.GetType()](update, callback)
	} else {
		gomer.SendMessage(chatID, i18n.Localize("oops"))
	}
}

func (state *ForcePublishState) publishConfirm(update *tgbotapi.Update, callback callback.Callback) {
	gomer := state.server.Gomer
	chatID := update.FromChat().ID

	gomer.SendMessage(chatID, i18n.Localize("forcePublishProgress"))

	buildDate := state.data.SelectedDate
	publishDate := &database.Date{Time: time.Now()}

	if err := state.server.ReportPublisher.Publish(buildDate, publishDate); err == nil {
		gomer.SendMessage(chatID, i18n.Localize("forcePublishSuccess"))
	} else {
		gomer.SendMessage(chatID, i18n.Localize("forcePublishError"))
	}

	state.stateMachine.
		Set(Idle, chatID, &StateContext{}).
		Init(update)
}

func (state *ForcePublishState) publishCancel(update *tgbotapi.Update, callback callback.Callback) {
	chatID := update.FromChat().ID

	state.stateMachine.
		Set(Idle, chatID, &StateContext{}).
		Init(update)
}
