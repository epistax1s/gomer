package strategy

import (
	"fmt"
	"strings"

	"github.com/epistax1s/gomer/internal/database"
	"github.com/epistax1s/gomer/internal/i18n"
	"github.com/epistax1s/gomer/internal/log"
	"github.com/epistax1s/gomer/internal/model"
	"github.com/epistax1s/gomer/internal/redmine"
	"github.com/epistax1s/gomer/internal/report/utils"
)

type RedmineExtStrategy struct {
	client          *redmine.RedmineClient
	excludePatterns []string
}

func NewRedmineExtStrategy(cl *redmine.RedmineClient, exc []string) *RedmineExtStrategy {
	return &RedmineExtStrategy{
		client:          cl,
		excludePatterns: exc,
	}
}

func (s *RedmineExtStrategy) FetchCommit(user *model.User, date *database.Date) string {
	entries, err := s.client.GetTimeEntries(int(*user.RedmineID), date.String())
	if err != nil {
		log.Error(
			"Error when uploading a commit from redmine",
			"user.ID", user.ID, "username", user.Username, "date", date)

		return "- " + i18n.Localize("commitRedmineConnectionError")
	}

	if len(entries) == 0 {
		log.Info(
			"No redmine entries",
			"user.ID", user.ID, "username", user.Username, "date", date)

		return "- " + i18n.Localize("commitDidntSent")
	}

	// Grouping by issues
	grouped := make(map[int]struct {
		subject  string
		comments []string
	})

	for _, entry := range entries {
		if _, ok := grouped[entry.Issue.ID]; !ok {
			subject, err := s.client.GetIssue(entry.Issue.ID)
			if err != nil {
				return "-" + i18n.Localize("commitRedmineConnectionError")
			}
			grouped[entry.Issue.ID] = struct {
				subject  string
				comments []string
			}{
				subject: subject,
			}
		}
		group := grouped[entry.Issue.ID]
		comments := strings.TrimSpace(entry.Comments)
		if !utils.ShouldExcludeComment(comments, s.excludePatterns) {
			group.comments = append(group.comments, comments)
		}
		grouped[entry.Issue.ID] = group
	}

	var builder strings.Builder
	for issueID, group := range grouped {
		// Skip issues with no comments after filtering
		if len(group.comments) == 0 {
			continue
		}

		// Format: #номерЗадачи Заголовок Задачи
		builder.WriteString(fmt.Sprintf("#%d %s\n", issueID, group.subject))

		// Add comments with "- " prefix
		for _, comment := range group.comments {
			builder.WriteString(fmt.Sprintf("- %s\n", comment))
		}
	}
	return builder.String()
}
