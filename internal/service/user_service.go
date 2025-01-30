package service

import (
	"errors"

	"gorm.io/gorm"

	"github.com/epistax1s/gomer/internal/log"
	"github.com/epistax1s/gomer/internal/model"
	"github.com/epistax1s/gomer/internal/repository"
)

type UserService interface {
	FindUserByChatID(int64) (*model.User, error)
	FindPaginated(page int, pageSize int) ([]model.User, error)
	FindAll() ([]model.User, error)
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

func (service *userService) FindUserByChatID(chatID int64) (*model.User, error) {
	user, err := service.userRepo.FindByChatID(chatID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Info(
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

func (service *userService) FindPaginated(page int, pageSize int) ([]model.User, error) {
	return service.userRepo.FindPaginated(page, pageSize)
}

func (service *userService) FindAll() ([]model.User, error) {
	return service.userRepo.FindAll()
}

func (service *userService) TrackUser(user *model.User) error {
	existsUser, _ := service.FindUserByChatID(user.ChatID)
	if existsUser != nil {
		// modifying and activating a previously deleted user
		existsUser.DepartmentId = user.DepartmentId
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
	user, err := service.FindUserByChatID(chatID)
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
