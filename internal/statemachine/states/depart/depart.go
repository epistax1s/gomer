package depart

import (
	"github.com/epistax1s/gomer/internal/server"
	. "github.com/epistax1s/gomer/internal/statemachine/core"
)

type TrackDepartmentState struct {
	server       *server.Server
	stateMachine *StateMachine
	data         *StateContext
}

func NewTrackDepartmentState(server *server.Server, stateMachine *StateMachine, data *StateContext) State {
	return &TrackDepartmentState{
		server:       server,
		stateMachine: stateMachine,
		data:         data,
	}
}
