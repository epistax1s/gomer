package strategy

import (
	"github.com/epistax1s/gomer/internal/database"
	"github.com/epistax1s/gomer/internal/model"
)

type CommitSourceStrategy interface {
    FetchCommit(user *model.User, date *database.Date) string
}
