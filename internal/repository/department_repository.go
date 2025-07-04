package repository

import (
	"fmt"

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
	result := repo.db.
		Debug().
		Where(fmt.Sprintf("%s != ?", model.DepartmentTypeColumn), model.DepartmentTypeSystem).
		First(&departments)

	return departments, result.Error
}
