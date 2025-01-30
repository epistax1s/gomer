package manage_departs

import (
	"github.com/epistax1s/gomer/internal/callback"
	"github.com/epistax1s/gomer/internal/server"
	. "github.com/epistax1s/gomer/internal/statemachine/core"
)

type ManageDepartsState struct {
	server       *server.Server
	stateMachine *StateMachine
	handlers     map[string]StateHandler
}

func NewManagaDepartsState(server *server.Server, stateMachine *StateMachine, data *StateContext) State {
	state := &ManageDepartsState{
		server:       server,
		stateMachine: stateMachine,
	}

	state.handlers = map[string]StateHandler{
		callback.PREV:         state.prevHandler,
		callback.NEXT:         state.nextHandler,
		callback.SWAP:         state.swapHandler,
		callback.ADD:          state.addHandeler,
		callback.SELECT:       state.selectHandler,
		callback.DELETE:       state.deleteHandler,
		callback.BACK_TO_LIST: state.listHandler,
		callback.EXIT:         state.exitHandler,
	}

	return state
}
