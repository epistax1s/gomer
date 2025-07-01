package strategy

import (
	"fmt"
	"strings"

	"github.com/epistax1s/gomer/internal/database"
	"github.com/epistax1s/gomer/internal/log"
	"github.com/epistax1s/gomer/internal/model"
	"github.com/epistax1s/gomer/internal/redmine"
)

type RedmineExtStrategy struct {
	client *redmine.RedmineClient
}

func NewRedmineExtStrategy(client *redmine.RedmineClient) *RedmineExtStrategy {
	return &RedmineExtStrategy{client: client}
}

func (s *RedmineExtStrategy) FetchCommit(user *model.User, date *database.Date) (string, bool) {
	entries, err := s.client.GetTimeEntries(int(user.ID), date.String())
	if err != nil {
		log.Error(
			"Error when uploading a commit from redmine",
			"user.ID", user.ID, "username", user.Username, "date", date)

		return "", false
	}

	if len(entries) == 0 {
		log.Info(
			"No redmine entries",
			"user.ID", user.ID, "username", user.Username, "date", date)

		return "", false
	}

	// Grouping by issues
	grouped := make(map[int]struct {
		subject  string
		hours    float64
		comments []string
	})

	for _, entry := range entries {
		if _, ok := grouped[entry.Issue.ID]; !ok {
			subject, err := s.client.GetIssue(entry.Issue.ID)
			if err != nil {
				return "", false
			}
			grouped[entry.Issue.ID] = struct {
				subject  string
				hours    float64
				comments []string
			}{
				subject: subject,
			}
		}
		group := grouped[entry.Issue.ID]
		group.hours += entry.Hours
		if entry.Comments != "" {
			group.comments = append(group.comments, entry.Comments)
		}
		grouped[entry.Issue.ID] = group
	}

	var builder strings.Builder
	for _, group := range grouped {
		builder.WriteString(fmt.Sprintf("- %s: %.2f hours", group.subject, group.hours))
		if len(group.comments) > 0 {
			builder.WriteString(" (")
			builder.WriteString(strings.Join(group.comments, "; "))
			builder.WriteString(")")
		}
		builder.WriteString("\n")
	}
	return builder.String(), true
}
