package repository

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/epistax1s/gomer/internal/model"
)

type AuthKeyRepository interface {
	FindByKey(string) (*model.AuthKey, error)
}

type authKeyRepository struct {
	db *gorm.DB
}

func NewAuthKeyRepository(db *gorm.DB) AuthKeyRepository {
	return &authKeyRepository{
		db: db,
	}
}

func (repo *authKeyRepository) FindByKey(key string) (*model.AuthKey, error) {
	var authKey model.AuthKey

	result := repo.db.
		Where(fmt.Sprintf("%s = ?", model.AuthKeyColumn), key).
		First(&authKey)

	return &authKey, result.Error
}
