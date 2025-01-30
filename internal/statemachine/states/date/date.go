package date

import (
	"github.com/epistax1s/gomer/internal/server"
	. "github.com/epistax1s/gomer/internal/statemachine/core"

	callback "github.com/epistax1s/gomer/internal/callback"
)

type DateState struct {
	server       *server.Server
	stateMachine *StateMachine
	data         *StateContext
	handlers     map[string]StateHandler
}

func NewDateState(server *server.Server, stateMachine *StateMachine, data *StateContext) State {
	state := &DateState{
		server:       server,
		stateMachine: stateMachine,
		data:         data,
	}

	state.handlers = map[string]StateHandler{
		callback.CALENDAR_PREV: state.prevHandler,
		callback.CALENDAR_NEXT: state.nextHandler,
		callback.CALENDAR_DATE: state.dateHandler,
	}

	return state
}
