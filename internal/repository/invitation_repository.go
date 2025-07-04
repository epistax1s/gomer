package repository

import (
	"gorm.io/gorm"

	"github.com/epistax1s/gomer/internal/model"
)

type InvitationRepository interface {
	Create(invite *model.Invitation) error
	FindByCode(code string) (*model.Invitation, error)
	Update(invitation *model.Invitation) error
	FindByCreatedBy(userID int64) ([]model.Invitation, error)
}

type invitationRepository struct {
	db *gorm.DB
}

func NewInvitationRepository(db *gorm.DB) InvitationRepository {
	return &invitationRepository{
		db: db,
	}
}

func (repo *invitationRepository) Create(invitation *model.Invitation) error {
	return repo.db.Debug().Create(invitation).Error
}

func (repo *invitationRepository) FindByCode(code string) (*model.Invitation, error) {
	var invitation model.Invitation
	result := repo.db.Where("code = ?", code).First(&invitation)
	return &invitation, result.Error
}

func (repo *invitationRepository) Update(invitation *model.Invitation) error {
	return repo.db.Updates(invitation).Error
}

func (repo *invitationRepository) FindByCreatedBy(userID int64) ([]model.Invitation, error) {
	var invitations []model.Invitation

	result := repo.db.
		Preload("CreatedBy").
		Preload("UsedBy").
		Where("created_by_user_id = ?", userID).
		Find(&invitations)

	return invitations, result.Error
}
