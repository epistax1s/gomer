package report

import (
	"time"

	"github.com/epistax1s/gomer/internal/log"
	"github.com/epistax1s/gomer/internal/server"
	"github.com/robfig/cron/v3"
)

func StartNotification(server *server.Server) {
	c := cron.New()
	c.AddFunc(server.Config.Report.NotificationCron, func() {
		now := time.Now()
		if !isBusinessDay(&now) {
			log.Info(
				"today is not a business day. The notification will not be sent",
				"action", "report notification cron")

			return
		}

		reportDate := calcReportDate(&now)

		log.Info(
			"sending a notification before generating a report",
			"action", "report notification cron", "reportDate", reportDate)

		publishNotification(server, reportDate)
	})
	c.Start()
}
