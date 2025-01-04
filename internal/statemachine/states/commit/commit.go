package commit

import (
	"github.com/epistax1s/gomer/internal/server"
	. "github.com/epistax1s/gomer/internal/statemachine/core"
)

type CommitState struct {
	server       *server.Server
	stateMachine *StateMachine
	data         *StateContext
}

func NewCommitState(server *server.Server, stateMachine *StateMachine, data *StateContext) State {
	return &CommitState{
		server:       server,
		stateMachine: stateMachine,
		data:         data,
	}
}
