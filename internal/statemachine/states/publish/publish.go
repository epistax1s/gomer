package publish

import (
	"github.com/epistax1s/gomer/internal/server"
	. "github.com/epistax1s/gomer/internal/statemachine/core"
)

type ForcePublishState struct {
	server       *server.Server
	stateMachine *StateMachine
	data         *StateContext
}

func NewForcePublishState(server *server.Server, stateMachine *StateMachine, data *StateContext) State {
	return &ForcePublishState{
		server:       server,
		stateMachine: stateMachine,
		data:         data,
	}
}
