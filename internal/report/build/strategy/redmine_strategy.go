package strategy

import (
	"strings"

	"github.com/epistax1s/gomer/internal/database"
	"github.com/epistax1s/gomer/internal/i18n"
	"github.com/epistax1s/gomer/internal/log"
	"github.com/epistax1s/gomer/internal/model"
	"github.com/epistax1s/gomer/internal/redmine"
	"github.com/epistax1s/gomer/internal/report/utils"
)

type RedmineStrategy struct {
	client          *redmine.RedmineClient
	excludePatterns []string
}

func NewRedmineStrategy(cl *redmine.RedmineClient, exc []string) *RedmineStrategy {
	return &RedmineStrategy{
		client:          cl,
		excludePatterns: exc,
	}
}

func (s *RedmineStrategy) FetchCommit(user *model.User, date *database.Date) string {
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

	var builder strings.Builder
	for i, entry := range entries {
		// skip comments that match regular expressions for exceptions.
		comments := strings.TrimSpace(entry.Comments)
		if utils.ShouldExcludeComment(comments, s.excludePatterns) {
			continue
		}

		builder.WriteString("- ")
		builder.WriteString(comments)
		if i < len(entries) {
			builder.WriteString("\n")
		}
	}

	return builder.String()
}
