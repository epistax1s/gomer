package publish

import (
	gomer "github.com/epistax1s/gomer/internal/bot"

	"time"

	"github.com/epistax1s/gomer/internal/config"
	"github.com/epistax1s/gomer/internal/database"
	"github.com/epistax1s/gomer/internal/log"
	"github.com/epistax1s/gomer/internal/report/build"
	"github.com/epistax1s/gomer/internal/report/utils"
	"github.com/epistax1s/gomer/internal/service"
	"github.com/robfig/cron/v3"
)

type ReportPublisher struct {
	groupService  service.GroupService
	reportConfig  *config.ReportConfig
	reportBuilder *build.ReportBuilder
	gomer         *gomer.Gomer
	cron          *cron.Cron
}

func NewReportPublisher(
	gs service.GroupService,
	rc *config.ReportConfig,
	rb *build.ReportBuilder,
	gr *gomer.Gomer,
) *ReportPublisher {

	return &ReportPublisher{
		groupService:  gs,
		reportConfig:  rc,
		reportBuilder: rb,
		gomer:         gr,
		cron:          cron.New(),
	}
}

// Start запускает cron-задачу, которая по расписанию вызывает генерацию отчёта.
func (rp *ReportPublisher) StartCron() {
	rp.cron.AddFunc(rp.reportConfig.PublishCron, func() {
		now := time.Now()
		if !utils.IsBusinessDay(&now) {
			log.Info(
				"today is not a business day. The report will not be generated",
				"action", "report publish cron")
			return
		}

		buildDate := utils.CastToBuildDate(&now)
		publishDate := &database.Date{Time: now}

		log.Info(
			"generating a report",
			"action", "report cron", "buildDate", buildDate, "publishDate", publishDate)

		if err := rp.Publish(buildDate, publishDate); err != nil {
			log.Error(
				"error generating the report",
				"action", "report cron", "date", buildDate, "err", err)
			return
		}

		log.Info(
			"the report has been successfully generated",
			"action", "report cron", "date", buildDate)
	})

	rp.cron.Start()
}

func (rp *ReportPublisher) Publish(buildDate *database.Date, publishDate *database.Date) error {
	messages, err := rp.reportBuilder.BuildDailyReport(buildDate, publishDate)

	groups, err := rp.groupService.FindAll()
	if err != nil {
		return err
	}

	for _, group := range groups {
		for _, message := range messages {
			rp.gomer.SendMessage(group.GroupID, message)
		}
	}

	return nil
}
