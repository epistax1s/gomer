package service

import (
	"errors"

	"github.com/epistax1s/gomer/internal/log"
	"github.com/epistax1s/gomer/internal/model"
	"github.com/epistax1s/gomer/internal/repository"
	"gorm.io/gorm"
)

type GroupService interface {
	LinkGroup(groupID int64, title string) error
	FindByID(id int64) (*model.Group, error)
	FindByGroupID(int64) (*model.Group, error)
	FindAll() ([]model.Group, error)
	FindPaginated(page int, pageSize int) ([]model.Group, error)
	CountAll() (int64, error)
}

type groupService struct {
	groupRepo repository.GroupRepository
}

func NewGroupService(groupRepo repository.GroupRepository) GroupService {
	return &groupService{
		groupRepo: groupRepo,
	}
}

func (service *groupService) LinkGroup(groupID int64, title string) error {
	if group, _ := service.FindByGroupID(groupID); group != nil {
		log.Info(
			"The group is already linked to the bot",
			"groupID", groupID, "title", title)

		return nil
	}

	return service.groupRepo.Create(
		&model.Group{
			GroupID: groupID,
			Title:   title,
		},
	)
}

func (service *groupService) FindByID(id int64) (*model.Group, error) {
	group, err := service.groupRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		log.Error(
			"error when searching for group by id",
			"id", id, "err", err)

		return nil, err
	}

	return group, nil
}

func (service *groupService) FindByGroupID(groupID int64) (*model.Group, error) {
	group, err := service.groupRepo.FindByGroupID(groupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		log.Error(
			"error when searching for for group by group id",
			"groupID", groupID, "err", err)

		return nil, err
	}

	return group, nil
}

func (service *groupService) FindAll() ([]model.Group, error) {
	return service.groupRepo.FindAll()
}

func (service *groupService) FindPaginated(page int, pageSize int) ([]model.Group, error) {
	return service.groupRepo.FindPaginated(page, pageSize)
}

func (service *groupService) CountAll() (int64, error) {
	return service.groupRepo.CountAll()
}
