package manage_groups

import (
	"github.com/epistax1s/gomer/internal/server"
	. "github.com/epistax1s/gomer/internal/statemachine/core"

	callback "github.com/epistax1s/gomer/internal/callback"
)

type ManageGroupsState struct {
	server       *server.Server
	stateMachine *StateMachine
	handlers     map[string]StateHandler
}

func NewManageGroupsState(server *server.Server, stateMachine *StateMachine, data *StateContext) State {
	state := &ManageGroupsState{
		server:       server,
		stateMachine: stateMachine,
	}

	state.handlers = map[string]StateHandler{
		callback.AG_PREV:       state.prevHandler,
		callback.AG_NEXT:       state.nextHandler,
		callback.AG_SELECT:     state.selectHandler,
		callback.AG_UNLINK:     state.unlinkHandler,
		callback.AG_GROUP_LIST: state.listHandler,
		callback.EXIT:          state.exitHandler,
	}

	return state
}
