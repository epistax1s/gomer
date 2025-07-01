package config

import (
	"github.com/epistax1s/gomer/internal/server"
	. "github.com/epistax1s/gomer/internal/statemachine/core"
)

const (
	cmdManualSrc     = "manual"
	cmdRedmineSrc    = "redmine"
	cmdRedmineExtSrc = "redmine_ext"
)

type ConfigState struct {
	server       *server.Server
	stateMachine *StateMachine
	data         *StateContext
	handlers     map[string]StateHandler
}

func NewConfigState(server *server.Server, stateMachine *StateMachine, data *StateContext) State {
	state := &ConfigState{
		server:       server,
		stateMachine: stateMachine,
		data:         data,
	}

	state.handlers = map[string]StateHandler{
		cmdManualSrc:     state.manualHandler,
		cmdRedmineSrc:    state.redmineHandler,
		cmdRedmineExtSrc: state.redmineExtHandler,
	}

	return state
}
