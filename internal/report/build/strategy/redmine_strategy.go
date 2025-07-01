package strategy

import (
	"fmt"
	"strings"

	"github.com/epistax1s/gomer/internal/database"
	"github.com/epistax1s/gomer/internal/log"
	"github.com/epistax1s/gomer/internal/model"
	"github.com/epistax1s/gomer/internal/redmine"
)

type RedmineStrategy struct {
	client *redmine.RedmineClient
}

func NewRedmineStrategy(client *redmine.RedmineClient) *RedmineStrategy {
	return &RedmineStrategy{client: client}
}

func (s *RedmineStrategy) FetchCommit(user *model.User, date *database.Date) (string, bool) {
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

	var builder strings.Builder
	for _, entry := range entries {
		builder.WriteString(fmt.Sprintf("- %s: %.2f hours (%s)\n",
			entry.Project.Name, entry.Hours, entry.Comments))
	}

	return builder.String(), true
}
