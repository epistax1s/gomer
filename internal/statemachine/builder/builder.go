package builder

import (
	"github.com/epistax1s/gomer/internal/server"
	"github.com/epistax1s/gomer/internal/statemachine/states/commit"
	"github.com/epistax1s/gomer/internal/statemachine/states/config"
	"github.com/epistax1s/gomer/internal/statemachine/states/date"
	"github.com/epistax1s/gomer/internal/statemachine/states/depart"
	"github.com/epistax1s/gomer/internal/statemachine/states/idle"
	"github.com/epistax1s/gomer/internal/statemachine/states/name"
	"github.com/epistax1s/gomer/internal/statemachine/states/publish"

	commit_modify "github.com/epistax1s/gomer/internal/statemachine/states/commit_modify"
	manage_groups "github.com/epistax1s/gomer/internal/statemachine/states/manage_groups"
	manage_users "github.com/epistax1s/gomer/internal/statemachine/states/manage_users"
	redmine_id "github.com/epistax1s/gomer/internal/statemachine/states/redmin_id"

	. "github.com/epistax1s/gomer/internal/statemachine/core"
)

func NewStateMachine(server *server.Server) *StateMachine {
	return &StateMachine{
		Server:       server,
		CurrentState: make(map[int64]State),
		StateFactory: map[StateType]StateFactory{
			Idle:            idle.NewIdleState,
			TrackDepartment: depart.NewTrackDepartmentState,
			TrackName:       name.NewTrackNameState,
			Date:            date.NewDateState,
			Commit:          commit.NewCommitState,
			CommitModify:    commit_modify.NewCommitModifyState,
			Config:          config.NewConfigState,
			ForcePublish:    publish.NewForcePublishState,
			RedmineID:       redmine_id.NewRedmineIDState,
			ManageUsers:     manage_users.NewManageUsersState,
			ManageGrops:     manage_groups.NewManageGroupsState,
		},
	}

}
