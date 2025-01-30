package publish

import (
	"time"

	"github.com/epistax1s/gomer/internal/database"
	"github.com/epistax1s/gomer/internal/i18n"
	"github.com/epistax1s/gomer/internal/log"
	"github.com/epistax1s/gomer/internal/report"
	. "github.com/epistax1s/gomer/internal/statemachine/core"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (state *ForcePublishState) Init(update *tgbotapi.Update) {
	state.Handle(update)
}

func (state *ForcePublishState) Handle(update *tgbotapi.Update) {
	gomer := state.server.Gomer
	chatID := update.FromChat().ID

	gomer.SendMessage(chatID, i18n.Localize("forcePublishProcess"))

	// TODO rename to SelectedDate
	reportDate := state.data.CommitDate
	publishDate := &database.Date{
		Time: time.Now(),
	}

	if err := report.BuildDailyReport(state.server, reportDate, publishDate); err == nil {
		log.Info(
			"Report was published successfully",
			"reportDate", reportDate, "publishDate", publishDate)

		gomer.SendMessage(chatID, i18n.Localize("forcePublishSuccess"))
	} else {
		log.Error(
			"Errors during the force publication of the report",
			"reportDate", reportDate, "publishDate", publishDate)

		gomer.SendMessage(chatID, i18n.Localize("forcePublishError"))
	}

	state.stateMachine.
		Set(Idle, chatID, &StateContext{}).
		Init(update)
}
