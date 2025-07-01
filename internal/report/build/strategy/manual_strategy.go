package strategy

import (
	"github.com/epistax1s/gomer/internal/database"
	"github.com/epistax1s/gomer/internal/log"
	"github.com/epistax1s/gomer/internal/model"
	"github.com/epistax1s/gomer/internal/service"
)

// ManualStrategy - данные из БД
type ManualStrategy struct {
	cm service.CommitService
}

func NewManualStrategy(cm service.CommitService) *ManualStrategy {
	return &ManualStrategy{cm: cm}
}

func (s *ManualStrategy) FetchCommit(user *model.User, date *database.Date) (string, bool) {
	commit, _ := s.cm.FindCommitByUserIdAndDate(user.ID, date)

	if commit == nil {
		log.Info(
			"Commit not found",
			"userID", user.ID, "username", user.Username, "date", date)

		return "", false
	}

	return commit.Payload, true
}
