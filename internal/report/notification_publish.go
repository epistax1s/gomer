package report

import (
	"fmt"

	"github.com/epistax1s/gomer/internal/database"
	"github.com/epistax1s/gomer/internal/i18n"
	"github.com/epistax1s/gomer/internal/server"
)

func publishNotification(server *server.Server, reportDate *database.Date) error {
	fullCommits, err := server.FullCommitService.FindAllByDate(reportDate)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf(i18n.Localize("notification"), reportDate.String())

	for _, fullCommit := range fullCommits {
		if !fullCommit.CommitSent {
			server.Gomer.SendMessage(fullCommit.ChatID, msg)
		}
	}

	return nil
}
