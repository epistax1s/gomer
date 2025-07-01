package build

import (
	"fmt"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/epistax1s/gomer/internal/database"
	"github.com/epistax1s/gomer/internal/i18n"
	"github.com/epistax1s/gomer/internal/model"
	"github.com/epistax1s/gomer/internal/redmine"
	"github.com/epistax1s/gomer/internal/report/build/strategy"
	"github.com/epistax1s/gomer/internal/service"
)

type ReportBuilder struct {
	userService    service.UserService
	strategies     map[string]strategy.CommitSourceStrategy
	messageMaxSize int64
}

func NewReportBuilder(
	us service.UserService,
	cs service.CommitService,
	rc *redmine.RedmineClient,
	messageMaxSize int64,
) *ReportBuilder {

	return &ReportBuilder{
		userService: us,
		strategies: map[string]strategy.CommitSourceStrategy{
			model.UserCommitSrcManual:     strategy.NewManualStrategy(cs),
			model.UserCommitSrcRedmine:    strategy.NewRedmineStrategy(rc),
			model.UserCommitSrcRedmineExt: strategy.NewRedmineExtStrategy(rc),
		},
	}
}

func (rb *ReportBuilder) BuildDailyReport(buildDate *database.Date, publishDate *database.Date) ([]string, error) {
	users, err := rb.userService.FindAllActive()
	if err != nil {
		return nil, err
	}

	sort.Sort(model.ByDepartmentOrderAndName(users))

	var messages []string
	var messageBuilder strings.Builder
	var departID int64 = -1

	messageBuilder.WriteString("#")
	messageBuilder.WriteString(i18n.Localize("reportTitle"))
	messageBuilder.WriteString(" ")
	messageBuilder.WriteString(publishDate.String())
	messageBuilder.WriteString("\n\n")

	for _, user := range users {
		var segmentBuilder strings.Builder
		department := user.Department

		if user.DepartmentID != departID {
			segmentBuilder.WriteString(fmt.Sprintf("——— %v ———\n\n", department.Name))
			departID = user.DepartmentID
		}

		strategy, ok := rb.strategies[user.CommitSrc]
		if !ok {
			strategy = rb.strategies[model.UserCommitSrcManual] // default to manual
		}

		commitPayload, commitExists := strategy.FetchCommit(&user, buildDate)
		if !commitExists {
			commitPayload = "- " + i18n.Localize("commitDidntSent")
			continue
		}

		segmentBuilder.WriteString(user.Name)
		segmentBuilder.WriteString(" @")
		segmentBuilder.WriteString(user.Username)
		segmentBuilder.WriteString(":\n")
		segmentBuilder.WriteString(commitPayload)
		segmentBuilder.WriteString("\n\n")

		if rb.isMessageFilled(&messageBuilder, &segmentBuilder) {
			messages = append(messages, messageBuilder.String())
			messageBuilder.Reset()
		}
		messageBuilder.WriteString(segmentBuilder.String())
	}

	messages = append(messages, messageBuilder.String())

	return messages, nil
}

func (rb *ReportBuilder) isMessageFilled(messageBuilder *strings.Builder, segmentBuilder *strings.Builder) bool {
	messageLen := utf8.RuneCountInString(messageBuilder.String())
	segmentLen := utf8.RuneCountInString(segmentBuilder.String())

	return messageLen+segmentLen > int(rb.messageMaxSize)
}
