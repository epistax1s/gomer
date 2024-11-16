package service

import (
	"errors"

	"gorm.io/gorm"

	"github.com/epistax1s/gomer/internal/model"
	"github.com/epistax1s/gomer/internal/repository"
	"github.com/epistax1s/gomer/internal/log"
)

type DepartService interface {
	FindById(int64) (*model.Department, error)
	FindAll() ([]model.Department, error)
}

type departService struct {
	departRepo repository.DepartRepository
}

func NewDepartService(departRepo repository.DepartRepository) DepartService {
	return &departService{
		departRepo: departRepo,
	}
}

func (service *departService) FindById(id int64) (*model.Department, error) {
	depart, err := service.departRepo.FindById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		log.Error(
			"error when searching for department by id",
			"id", id, "err", err)

		return nil, err
	}

	return depart, nil
}

func (service *departService) FindAll() ([]model.Department, error) {
	return service.departRepo.FindAll()
}
