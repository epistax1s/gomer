package core

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/epistax1s/gomer/internal/callback"
	"github.com/epistax1s/gomer/internal/database"
	"github.com/epistax1s/gomer/internal/model"
	"github.com/epistax1s/gomer/internal/server"
)

// заменить StateMachine на интерфейс и мб разбить на файлы
type State interface {
	Init(*tgbotapi.Update)
	Handle(*tgbotapi.Update)
}

type StateFactory func(*server.Server, *StateMachine, *StateContext) State

type StateHandler func(*tgbotapi.Update, callback.Callback)

type StateMachine struct {
	Server       *server.Server
	CurrentState map[int64]State
	StateFactory map[StateType]StateFactory
}

type StateType string

const (
	Idle            StateType = "idle"
	TrackDepartment StateType = "trackDepartment"
	TrackName       StateType = "trackName"
	Date            StateType = "date"
	Commit          StateType = "commit"
	CommitModify    StateType = "commitModify"
	Config			StateType = "config"
	ForcePublish    StateType = "Publish"
	ManageUsers     StateType = "manageUsers"
	ManageGrops     StateType = "manageGroups"
)

type StateContext struct {
	Commit     *model.Commit
	CommitDate *database.Date
	Department *model.Department
	NextState  StateType
}

func (stateMachine *StateMachine) Set(stateType StateType, chatID int64, data *StateContext) State {
	stateFactory, exists := stateMachine.StateFactory[stateType]
	if !exists {
		panic(
			fmt.Sprintf("no state factory was found for a state with type = %s", stateType),
		)
	}

	stateMachine.CurrentState[chatID] = stateFactory(stateMachine.Server, stateMachine, data)
	return stateMachine.CurrentState[chatID]
}

func (stateMachine *StateMachine) Get(chatID int64) State {
	state, exists := stateMachine.CurrentState[chatID]
	if !exists {
		state = stateMachine.Set(Idle, chatID, &StateContext{})
	}
	return state
}
