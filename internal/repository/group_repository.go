package repository

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/epistax1s/gomer/internal/model"
)

type GroupRepository interface {
	Create(*model.Group) error
	FindByID(int64) (*model.Group, error)
	FindByGroupID(int64) (*model.Group, error)
	FindAll() ([]model.Group, error)
	FindPaginated(page int, pageSize int) ([]model.Group, error)
	CountAll() (int64, error)
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

func (repo *groupRepository) FindByID(id int64) (*model.Group, error) {
	var group model.Group
	result := repo.db.First(&group, id)
	return &group, result.Error
}

func (repo *groupRepository) FindByGroupID(groupID int64) (*model.Group, error) {
	var group model.Group

	result := repo.db.
		Where(fmt.Sprintf("%s = ?", model.GroupGroupID), groupID).
		First(&group)

	return &group, result.Error
}

func (repo *groupRepository) FindAll() ([]model.Group, error) {
	var groups []model.Group
	result := repo.db.Find(&groups)
	return groups, result.Error
}

func (repo *groupRepository) FindPaginated(page int, pageSize int) ([]model.Group, error) {
	var groups []model.Group

	offset := (page - 1) * pageSize

	result := repo.db.Offset(offset).Limit(pageSize).Find(&groups)
	return groups, result.Error
}

func (repo *groupRepository) CountAll() (int64, error) {
	var count int64
	result := repo.db.Model(&model.Group{}).Count(&count)
	return count, result.Error
}
