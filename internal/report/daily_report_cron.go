package report

import (
	"time"

	"github.com/robfig/cron/v3"

	"github.com/epistax1s/gomer/internal/database"
	"github.com/epistax1s/gomer/internal/log"
	"github.com/epistax1s/gomer/internal/server"
)

func StartCron(server *server.Server) {
	c := cron.New(cron.WithSeconds())
	c.AddFunc(server.Config.Report.PublishCron, func() {
		now := time.Now()
		if !isBusinessDay(&now) {
			log.Info(
				"today is not a business day. The report will not be generated",
				"action", "report cron")

			return
		}

		reportDate := calcReportDate(&now)
		log.Info(
			"generating a report",
			"action", "report cron", "date", reportDate)

		err := BuildDailyReport(server, reportDate)
		if err != nil {
			log.Error(
				"error generating the report",
				"action", "report cron", "date", reportDate, "err", err)

			return
		}

		log.Info(
			"the report has been successfully generated",
			"action", "report cron", "date", reportDate)
	})
	c.Start()
}

func isBusinessDay(t *time.Time) bool {
	return t.Weekday() != time.Saturday && t.Weekday() != time.Sunday
}

func calcReportDate(t *time.Time) *database.Date {
	if t.Weekday() == time.Monday {
		reportTime := t.Add(-72 * time.Hour)
		return &database.Date{Time: reportTime}
	} else {
		reportTime := t.Add(-24 * time.Hour)
		return &database.Date{Time: reportTime}
	}
}