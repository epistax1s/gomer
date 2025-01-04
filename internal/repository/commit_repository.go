package repository

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/epistax1s/gomer/internal/database"
	"github.com/epistax1s/gomer/internal/model"
)

type CommitRepository interface {
	Create(*model.Commit) error
	Update(*model.Commit) error
	FindById(int64) (*model.Commit, error)
	FindByUserIdAndDate(int64, *database.Date) (*model.Commit, error)
	FindAllByDate(*database.Date) ([]model.Commit, error)
}

type commitRepository struct {
	db *gorm.DB
}

func NewCommitRepository(db *gorm.DB) CommitRepository {
	return &commitRepository{
		db: db,
	}
}

func (repo *commitRepository) Create(commit *model.Commit) error {
	return repo.db.Create(commit).Error
}

func (repo *commitRepository) Update(commit *model.Commit) error {
	return repo.db.Updates(commit).Error
}

func (repo *commitRepository) FindById(id int64) (*model.Commit, error) {
	var commit model.Commit
	result := repo.db.First(&commit, id)
	return &commit, result.Error
}

func (repo *commitRepository) FindByUserIdAndDate(userID int64, date *database.Date) (*model.Commit, error) {
	var commit model.Commit

	result := repo.db.
		Where(fmt.Sprintf("%s = ? AND %s = ?", model.CommitUserIDColumn, model.CommitDateColumn), userID, date).
		First(&commit)

	return &commit, result.Error
}

func (repo *commitRepository) FindAllByDate(date *database.Date) ([]model.Commit, error) {
	var commits []model.Commit

	result := repo.db.
		Preload(model.UserTable).
		Where(fmt.Sprintf("%s = ?", model.CommitDateColumn), date).
		Find(&commits)

	return commits, result.Error
}
