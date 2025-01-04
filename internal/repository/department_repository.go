package repository

import (
	"gorm.io/gorm"

	"github.com/epistax1s/gomer/internal/model"
)

type DepartRepository interface {
	FindById(int64) (*model.Department, error)
	FindAll() ([]model.Department, error)
}

type departRepository struct {
	db *gorm.DB
}

func NewDepartRepository(db *gorm.DB) DepartRepository {
	return &departRepository{
		db: db,
	}
}

func (repo *departRepository) FindById(id int64) (*model.Department, error) {
	var department model.Department
	result := repo.db.First(&department, id)
	return &department, result.Error
}

func (repo *departRepository) FindAll() ([]model.Department, error) {
	var departments []model.Department
	result := repo.db.Find(&departments)
	return departments, result.Error
}
