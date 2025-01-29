package state

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/epistax1s/gomer/internal/database"
	"github.com/epistax1s/gomer/internal/model"
	"github.com/epistax1s/gomer/internal/server"
)

type stateMachine struct {
	CurrentState map[int64]State
}

type State interface {
	Init(server *server.Server, update *tgbotapi.Update)
	Handle(server *server.Server, update *tgbotapi.Update)
}

type StateType string

const (
	Idle            StateType = "idle"
	TrackDepartment StateType = "trackDepartment"
	TrackName       StateType = "trackName"
	Date            StateType = "date"
	Commit          StateType = "commit"
	CommitModify    StateType = "commitModify"
	AdminGroup      StateType = "adminGroup"
)

type StateContext struct {
	Commit     *model.Commit
	CommitDate *database.Date
	Department *model.Department
	NextState  StateType
}

var statesFactory = map[StateType]func(data *StateContext) State{
	Idle:            NewIdleState,
	TrackDepartment: NewTrackDepartmentState,
	TrackName:       NewTrackNameState,
	Date:            NewDateState,
	Commit:          NewCommitState,
	CommitModify:    NewCommitModifyState,
	AdminGroup:      NewAdminGroupState,
}
