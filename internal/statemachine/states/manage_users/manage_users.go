package manage_users

import (
	"github.com/epistax1s/gomer/internal/server"
	. "github.com/epistax1s/gomer/internal/statemachine/core"

	callback "github.com/epistax1s/gomer/internal/callback"
)

type ManageUsersState struct {
	server       *server.Server
	stateMachine *StateMachine
	handlers     map[string]StateHandler
}

// TODO date -> context
func NewManageUsersState(server *server.Server, stateMachine *StateMachine, data *StateContext) State {
	state := &ManageUsersState{
		server:       server,
		stateMachine: stateMachine,
	}

	state.handlers = map[string]StateHandler{
		callback.PREV:   state.prevHandler,
		callback.NEXT:   state.nextHandler,
		callback.SELECT: state.selectHandler,
	}

	return state
}
