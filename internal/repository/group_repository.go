package repository

import (
	"gorm.io/gorm"

	"github.com/epistax1s/gomer/internal/model"
)

type GroupRepository interface {
	Create(*model.Group) error
	FindAll() ([]model.Group, error)
}

type groupRepository struct {
	db *gorm.DB
}

func NewGroupRepository(db *gorm.DB) GroupRepository {
	return &groupRepository{
		db: db,
	}
}

func (repo *groupRepository) Create(group *model.Group) error {
	return repo.db.Create(group).Error
}

func (repo *groupRepository) FindAll() ([]model.Group, error) {
	var groups []model.Group
	result := repo.db.Find(&groups)
	return groups, result.Error
}
