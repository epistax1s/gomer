package service

import (
	"errors"

	"gorm.io/gorm"

	"github.com/epistax1s/gomer/internal/log"
	"github.com/epistax1s/gomer/internal/model"
	"github.com/epistax1s/gomer/internal/repository"
)

type UserService interface {
	Create(*model.User) error
	Save(*model.User) error
	FindByID(int64) (*model.User, error)
	FindByChatID(int64) (*model.User, error)
	FindAll() ([]model.User, error)
	FindAllActive() ([]model.User, error)
	TrackUser(*model.User) error
	UntrackUser(int64) error
	UserExists(int64) (bool, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (service *userService) Create(user *model.User) error {
	return service.userRepo.Create(user)
}

func (service *userService) Save(user *model.User) error {
	return service.userRepo.Update(user)
}

func (service *userService) FindByID(id int64) (*model.User, error) {
	user, err := service.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Debug(
				"the user by id was not found",
				"id", id)

			return nil, nil
		}

		log.Error(
			"error when trying to find a user by id",
			"id", id, "err", err)

		return nil, err
	}

	return user, nil
}

func (service *userService) FindByChatID(chatID int64) (*model.User, error) {
	user, err := service.userRepo.FindByChatID(chatID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Debug(
				"the user by chat Id was not found",
				"chatID", chatID)

			return nil, nil
		}

		log.Error(
			"error when trying to find a user by chatID",
			"chatID", chatID, "err", err)

		return nil, err
	}

	return user, nil
}

func (service *userService) FindAll() ([]model.User, error) {
	return service.userRepo.FindAll()
}

func (service *userService) FindAllActive() ([]model.User, error) {
	return service.userRepo.FindAllActive()
}

func (service *userService) TrackUser(user *model.User) error {
	existsUser, _ := service.FindByChatID(user.ChatID)
	if existsUser != nil {
		// modifying and activating a previously deleted user
		existsUser.DepartmentID = user.DepartmentID
		existsUser.Name = user.Name
		existsUser.Username = user.Username
		existsUser.Status = model.UserStatusActive

		return service.userRepo.Update(existsUser)
	} else {
		// сreating a new user
		return service.userRepo.Create(user)
	}
}

func (service *userService) UntrackUser(chatID int64) error {
	user, err := service.FindByChatID(chatID)
	if err != nil {
		log.Error(
			"error when untracking a user",
			"chatID", chatID, "err", err)

		return err
	}

	user.Status = "deleted"

	if err := service.userRepo.Update(user); err != nil {
		log.Error(
			"error when untracking a user",
			"chatID", chatID, "err", err)

		return err
	}

	log.Info(
		"the user is no longer being tracked",
		"chatID", chatID)

	return nil
}

func (service *userService) UserExists(chatID int64) (bool, error) {
	user, err := service.userRepo.FindByChatID(chatID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return user.Status != model.UserStatusDeleted, nil
}

func (service *userService) CountAll() (int64, error) {
	return service.userRepo.CountAll()
}
