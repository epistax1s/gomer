package service

import (
	"errors"
	"testing"

	"github.com/epistax1s/gomer/internal/log/test"
	"github.com/epistax1s/gomer/internal/model"
	"github.com/epistax1s/gomer/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var errMockDB = errors.New("mock db error")

func init() {
	test.InitTestLogger()
}

func TestInvitationService_GenerateInvite(t *testing.T) {
	mockRepo := &mocks.MockInvitationRepository{}
	service := NewInvitationService(mockRepo)

	t.Run("successful invitation generation", func(t *testing.T) {
		userID := int64(1)
		mockRepo.On("Create", mock.AnythingOfType("*model.Invitation")).Return(nil).Once()

		code, err := service.GenerateInvite(userID)

		assert.NoError(t, err)
		assert.NotEmpty(t, code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		userID := int64(1)
		mockRepo.On("Create", mock.AnythingOfType("*model.Invitation")).Return(errMockDB).Once()

		code, err := service.GenerateInvite(userID)

		assert.Error(t, err)
		assert.Empty(t, code)
		mockRepo.AssertExpectations(t)
	})
}

func TestInvitationService_ValidateInvite(t *testing.T) {
	mockRepo := &mocks.MockInvitationRepository{}
	service := NewInvitationService(mockRepo)

	t.Run("valid active invitation", func(t *testing.T) {
		invitation := &model.Invitation{
			Code: "valid_code",
			Used: true,
		}
		mockRepo.On("FindByCode", "valid_code").Return(invitation, nil).Once()

		isValid := service.ValidateInvite("valid_code")

		assert.True(t, isValid)
		mockRepo.AssertExpectations(t)
	})

	t.Run("inactive invitation", func(t *testing.T) {
		invitation := &model.Invitation{
			Code: "used_code",
			Used: false,
		}
		mockRepo.On("FindByCode", "used_code").Return(invitation, nil).Once()

		isValid := service.ValidateInvite("used_code")

		assert.False(t, isValid)
		mockRepo.AssertExpectations(t)
	})

	t.Run("non-existent invitation", func(t *testing.T) {
		mockRepo.On("FindByCode", "invalid_code").Return(nil, errMockDB).Once()

		isValid := service.ValidateInvite("invalid_code")

		assert.False(t, isValid)
		mockRepo.AssertExpectations(t)
	})
}

func TestInvitationService_UseInvitation(t *testing.T) {
	mockRepo := &mocks.MockInvitationRepository{}
	service := NewInvitationService(mockRepo)

	t.Run("successful use of invitation", func(t *testing.T) {
		invitation := &model.Invitation{
			Code: "valid_code",
			Used: true,
		}
		user := &model.User{
			ID: 1,
		}

		mockRepo.On("FindByCode", "valid_code").Return(invitation, nil).Once()
		mockRepo.On("Update", mock.AnythingOfType("*model.Invitation")).Return(nil).Once()

		err := service.UseInvitation("valid_code", user)

		assert.NoError(t, err)
		assert.False(t, invitation.Used)
		assert.NotNil(t, invitation.UsedByID)
		assert.NotNil(t, invitation.UsedAt)
		mockRepo.AssertExpectations(t)
	})

	t.Run("already used invitation", func(t *testing.T) {
		invitation := &model.Invitation{
			Code: "used_code",
			Used: false,
		}
		user := &model.User{
			ID: 1,
		}

		mockRepo.On("FindByCode", "used_code").Return(invitation, nil).Once()

		err := service.UseInvitation("used_code", user)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("non-existent invitation", func(t *testing.T) {
		user := &model.User{
			ID: 1,
		}

		mockRepo.On("FindByCode", "invalid_code").Return(nil, errMockDB).Once()

		err := service.UseInvitation("invalid_code", user)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestInvitationService_GetInvitesByCreator(t *testing.T) {
	mockRepo := &mocks.MockInvitationRepository{}
	service := NewInvitationService(mockRepo)

	t.Run("successful retrieval", func(t *testing.T) {
		userID := int64(1)
		expectedInvites := []model.Invitation{
			{
				ID:          1,
				Code:        "code1",
				CreatedByID: userID,
				Used:        true,
			},
			{
				ID:          2,
				Code:        "code2",
				CreatedByID: userID,
				Used:        false,
			},
		}

		mockRepo.On("FindByCreatedBy", userID).Return(expectedInvites, nil).Once()

		invites, err := service.GetInvitesByCreator(userID)

		assert.NoError(t, err)
		assert.Equal(t, expectedInvites, invites)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		userID := int64(1)
		mockRepo.On("FindByCreatedBy", userID).Return([]model.Invitation{}, errMockDB).Once()

		invites, err := service.GetInvitesByCreator(userID)

		assert.Error(t, err)
		assert.Empty(t, invites)
		mockRepo.AssertExpectations(t)
	})
}
