package report

import (
	"fmt"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/epistax1s/gomer/internal/database"
	"github.com/epistax1s/gomer/internal/i18n"
	"github.com/epistax1s/gomer/internal/log"
	"github.com/epistax1s/gomer/internal/model"
	"github.com/epistax1s/gomer/internal/server"
)

func BuildDailyReport(server *server.Server, date *database.Date) error {
	reports, err := server.ReportService.FindAllByDate(date)
	if err != nil {
		log.Error(
			"errors in Compiling the daily report",
			"date", date, "err", err)
	}

	sort.Sort(model.ByDepartamentAndName(reports))

	var messages []string
	var messageBuilder strings.Builder
	var departID int64 = -1

	messageBuilder.WriteString("#")
	messageBuilder.WriteString(i18n.Localize("reportTitle"))
	messageBuilder.WriteString(" ")
	messageBuilder.WriteString(date.String())
	messageBuilder.WriteString("\n\n")

	for _, report := range reports {
		var segmentBuilder strings.Builder
		if report.DepartmentID != departID {
			segmentBuilder.WriteString(fmt.Sprintf("——————————— %v ———————————\n\n", report.DepartmentName))
			departID = report.DepartmentID
		}

		segmentBuilder.WriteString(buildDailyReportSegment(&report))
		if isMessageFilled(&messageBuilder, &segmentBuilder) {
			messages = append(messages, messageBuilder.String())
			messageBuilder.Reset()
		}

		messageBuilder.WriteString(segmentBuilder.String())
	}

	messages = append(messages, messageBuilder.String())

	if err := publishReport(server, messages); err != nil {
		log.Error(
			"errors when sending a report",
			"err", err)

		return err
	}

	return nil
}

func isMessageFilled(messageBuilder *strings.Builder, segmentBuilder *strings.Builder) bool {
	messageLen := utf8.RuneCountInString(messageBuilder.String())
	segmentLen := utf8.RuneCountInString(segmentBuilder.String())
	return messageLen+segmentLen > 4096
}

func buildDailyReportSegment(report *model.FullCommit) string {
	var reportBuilder strings.Builder

	reportBuilder.WriteString(report.Name)
	reportBuilder.WriteString(" @")
	reportBuilder.WriteString(report.Username)
	reportBuilder.WriteString(":\n")

	if !report.CommitSent {
		reportBuilder.WriteString("- ")
		reportBuilder.WriteString(i18n.Localize("commitDidntSent"))
	} else {
		reportBuilder.WriteString(report.CommitPayload)
	}

	reportBuilder.WriteString("\n\n")

	return reportBuilder.String()
}
