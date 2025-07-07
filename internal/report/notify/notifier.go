package notify

import (
	"fmt"
	"time"

	gomer "github.com/epistax1s/gomer/internal/bot"
	"github.com/epistax1s/gomer/internal/config"
	"github.com/epistax1s/gomer/internal/database"
	"github.com/epistax1s/gomer/internal/i18n"
	"github.com/epistax1s/gomer/internal/log"
	"github.com/epistax1s/gomer/internal/report/utils"
	"github.com/epistax1s/gomer/internal/service"
	"github.com/robfig/cron/v3"
)

type ReportNotifier struct {
	cron         *cron.Cron
	gomer        *gomer.Gomer
	reportConfig *config.ReportConfig
	userService  service.UserService
}

func NewReportNotifier(gr *gomer.Gomer, rc *config.ReportConfig, us service.UserService) *ReportNotifier {
	return &ReportNotifier{
		cron:         cron.New(),
		gomer:        gr,
		reportConfig: rc,
		userService:  us,
	}
}

func (n *ReportNotifier) StartCron() {
	n.cron.AddFunc(n.reportConfig.NotificationCron, func() {
		now := time.Now()
		if !utils.IsBusinessDay(&now) {
			log.Info(
				"today is not a business day. The notification will not be sent",
				"action", "report notification cron")
			return
		}

		buildDate := utils.CastToBuildDate(&now)

		log.Info(
			"sending a notification before generating a report",
			"action", "report notification cron", "buildDate", buildDate)

		if err := n.Notify(buildDate); err != nil {
			log.Error(
				"error sending notification",
				"action", "report notification cron", "err", err)
		}
	})

	n.cron.Start()
}

func (n *ReportNotifier) Notify(reportDate *database.Date) error {
	users, err := n.userService.FindAllActive()

	if err != nil {
		log.Error("Error when trying to send a notification to the user", "err", err)
		return err
	}

	for _, user := range users {
		n.gomer.SendMessage(user.ChatID, fmt.Sprintf(i18n.Localize("notification"), reportDate.String()))
	}

	return nil
}
