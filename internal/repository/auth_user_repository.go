package repository

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/epistax1s/gomer/internal/model"
)

type AuthUserRepository interface {
	Create(user *model.AuthUser) error
	FindByChatID(chatID int64) (*model.AuthUser, error)
}

type authUserRepository struct {
	db *gorm.DB
}

func NewAuthUserRepository(db *gorm.DB) AuthUserRepository {
	return &authUserRepository{
		db: db,
	}
}

func (repo *authUserRepository) Create(authUser *model.AuthUser) error {
	return repo.db.Create(authUser).Error
}

func (repo *authUserRepository) FindByChatID(chatID int64) (*model.AuthUser, error) {
	var authUser model.AuthUser

	result := repo.db.
		Where(fmt.Sprintf("%s = ?", model.AuthUserChatIDColumn), chatID).
		First(&authUser)

	return &authUser, result.Error
}
