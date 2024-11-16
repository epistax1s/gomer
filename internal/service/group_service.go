package service

import (
	"github.com/epistax1s/gomer/internal/model"
	"github.com/epistax1s/gomer/internal/repository"
)

type GroupService interface {
	LinkGroup(gropuID int64) error
	FindAll() ([]model.Group, error)
}

type groupService struct {
	groupRepo repository.GroupRepository
}

func NewGroupService(groupRepo repository.GroupRepository) GroupService {
	return &groupService{
		groupRepo: groupRepo,
	}
}

func (service *groupService) LinkGroup(groupID int64) error {
	return service.groupRepo.Create(
		&model.Group{
			GroupID: groupID,
		},
	)
}

func (service *groupService) FindAll() ([]model.Group, error) {
	return service.groupRepo.FindAll()
}
