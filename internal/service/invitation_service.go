package service

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/epistax1s/gomer/internal/log"
	"github.com/epistax1s/gomer/internal/model"
	"github.com/epistax1s/gomer/internal/repository"
)

type InvitationService interface {
	GenerateInvite(createdByUserID int64) (string, error)
	ValidateInvite(code string) bool
	UseInvitation(code string, usedByUser *model.User) error
	GetInvitesByCreator(userID int64) ([]model.Invitation, error)
}

type invitationService struct {
	invitationRepo repository.InvitationRepository
}

func NewInvitationService(invitationRepo repository.InvitationRepository) InvitationService {
	return &invitationService{
		invitationRepo: invitationRepo,
	}
}

func (service *invitationService) GenerateInvite(createdByUserID int64) (string, error) {
	// Generate a random code
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	code := base64.URLEncoding.EncodeToString(b)

	// Create invitation record
	invitation := &model.Invitation{
		Code:        code,
		CreatedByID: createdByUserID,
		CreatedAt:   time.Now().Format("2006-01-02 15:04:05"),
		IsActive:    true,
	}

	if err := service.invitationRepo.Create(invitation); err != nil {
		log.Error("failed to create invitation", "error", err)
		return "", fmt.Errorf("failed to create invitation: %w", err)
	}

	return code, nil
}

func (service *invitationService) ValidateInvite(code string) bool {
	// Find and validate the invitation
	invitation, err := service.invitationRepo.FindByCode(code)
	if err != nil {
		log.Info("failed to find invitation", "code", code)
		return false
	}

	// Validate invitation is still active
	if !invitation.IsActive {
		log.Info("invitation code has already been used", "code", code)
		return false
	}

	return true
}

func (service *invitationService) UseInvitation(code string, usedByUser *model.User) error {
	invitation, err := service.invitationRepo.FindByCode(code)
	if err != nil {
		log.Info("failed to find invitation", "code", code)
		return err
	}

	if !invitation.IsActive {
		return fmt.Errorf("invitation code %s has already been used", code)
	}

	// Mark invitation as used
	now := time.Now().Format("2006-01-02 15:04:05")
	invitation.UsedByID = &usedByUser.ID
	invitation.UsedAt = &now
	invitation.IsActive = false

	if err := service.invitationRepo.Update(invitation); err != nil {
		return fmt.Errorf("failed to mark invitation as used: %w", err)
	}

	return nil
}

func (service *invitationService) GetInvitesByCreator(userID int64) ([]model.Invitation, error) {
	return service.invitationRepo.FindByCreatedBy(userID)
}
