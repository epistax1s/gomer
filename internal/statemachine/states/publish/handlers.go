package publish

import (
	"github.com/epistax1s/gomer/internal/database"
	"github.com/epistax1s/gomer/internal/i18n"
	. "github.com/epistax1s/gomer/internal/statemachine/core"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

func (state *ForcePublishState) Init(update *tgbotapi.Update) {
	state.Handle(update)
}

func (state *ForcePublishState) Handle(update *tgbotapi.Update) {
	gomer := state.server.Gomer
	chatID := update.FromChat().ID

	gomer.SendMessage(chatID, i18n.Localize("forcePublishProcess"))

	buildDate := state.data.SelectedDate
	publishDate := &database.Date{Time: time.Now()}

	if err := state.server.ReportPublisher.Publish(buildDate, publishDate); err != nil {
		gomer.SendMessage(chatID, i18n.Localize("forcePublishSuccess"))
	} else {
		gomer.SendMessage(chatID, i18n.Localize("forcePublishError"))
	}

	state.stateMachine.
		Set(Idle, chatID, &StateContext{}).
		Init(update)
}
