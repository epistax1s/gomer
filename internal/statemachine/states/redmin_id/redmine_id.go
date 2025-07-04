package redminid


import (
	"github.com/epistax1s/gomer/internal/server"
	. "github.com/epistax1s/gomer/internal/statemachine/core"
)

type RedmineIDState struct {
	server       *server.Server
	stateMachine *StateMachine
	data         *StateContext
}

func NewRedmineIDState(server *server.Server, stateMachine *StateMachine, data *StateContext) State {
	return &RedmineIDState{
		server:       server,
		stateMachine: stateMachine,
		data:         data,
	}
}