package name

import (
	"github.com/epistax1s/gomer/internal/server"
	. "github.com/epistax1s/gomer/internal/statemachine/core"
)

type TrackNameState struct {
	server       *server.Server
	stateMachine *StateMachine
	data         *StateContext
}

func NewTrackNameState(server *server.Server, stateMachine *StateMachine, data *StateContext) State {
	return &TrackNameState{
		server:       server,
		stateMachine: stateMachine,
		data:         data,
	}
}
