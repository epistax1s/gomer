package idle

import (
	"github.com/epistax1s/gomer/internal/server"
	. "github.com/epistax1s/gomer/internal/statemachine/core"
)

const (
	cmdStart        = "start"
	cmdHelp         = "help"
	cmdTrack        = "track"
	cmdUntrack      = "untrack"
	cmdCommit       = "commit"
	cmdCommitModify = "modify"
	cmdForcePublish = "publish"
	cmdManageUsers  = "users"
	cmdManageGroups = "groups"
)

type IdleState struct {
	server       *server.Server
	stateMachine *StateMachine
	handlers     map[string]StateHandler
}

func NewIdleState(server *server.Server, stateMachine *StateMachine, data *StateContext) State {
	state := &IdleState{
		server:       server,
		stateMachine: stateMachine,
	}

	state.handlers = map[string]StateHandler{
		cmdStart:        state.helpHandler,
		cmdHelp:         state.helpHandler,
		cmdTrack:        state.trackHandler,
		cmdUntrack:      state.untrackHandler,
		cmdCommit:       state.commitHandler,
		cmdCommitModify: state.commitModifyHandler,
		cmdForcePublish: state.forcePublishHandler,
		cmdManageUsers:  state.manageUsersHandler,
		cmdManageGroups: state.manageGroupsHandler,
	}

	return state
}
