package publish

import (
	"github.com/epistax1s/gomer/internal/callback"
	"github.com/epistax1s/gomer/internal/server"
	. "github.com/epistax1s/gomer/internal/statemachine/core"
)

type ForcePublishState struct {
	server       *server.Server
	stateMachine *StateMachine
	data         *StateContext
	handlers     map[string]StateHandler
}

func NewForcePublishState(server *server.Server, stateMachine *StateMachine, data *StateContext) State {
	state := &ForcePublishState{
		server:       server,
		stateMachine: stateMachine,
		data:         data,
	}

	state.handlers = map[string]StateHandler{
		callback.ACTION_CONFIRM: state.publishConfirm,
		callback.ACTION_CANCEL:  state.publishCancel,
	}

	return state
}
