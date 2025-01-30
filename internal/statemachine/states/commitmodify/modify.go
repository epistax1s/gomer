package commit_modify

import (
	"github.com/epistax1s/gomer/internal/server"
	. "github.com/epistax1s/gomer/internal/statemachine/core"
)

type CommitModifyState struct {
	server       *server.Server
	stateMachine *StateMachine
	data         *StateContext
}

func NewCommitModifyState(server *server.Server, stateMachine *StateMachine, data *StateContext) State {
	return &CommitModifyState{
		server:       server,
		stateMachine: stateMachine,
		data:         data,
	}
}
