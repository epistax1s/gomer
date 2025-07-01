package mocks

import (
	"github.com/epistax1s/gomer/internal/model"
	"github.com/stretchr/testify/mock"
)

type MockInvitationRepository struct {
	mock.Mock
}

func (m *MockInvitationRepository) Create(invite *model.Invitation) error {
	args := m.Called(invite)
	return args.Error(0)
}

func (m *MockInvitationRepository) FindByCode(code string) (*model.Invitation, error) {
	args := m.Called(code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Invitation), args.Error(1)
}

func (m *MockInvitationRepository) Update(invitation *model.Invitation) error {
	args := m.Called(invitation)
	return args.Error(0)
}

func (m *MockInvitationRepository) FindByCreatedBy(userID int64) ([]model.Invitation, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.Invitation), args.Error(1)
}
