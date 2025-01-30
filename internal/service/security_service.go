package service

import (
	"errors"

	"github.com/epistax1s/gomer/internal/model"
	"gorm.io/gorm"
)

type SecurityService interface {
	IsAdmin(chatID int64) (bool, error)
}

type securityService struct {
	userService UserService
}

func NewSecurityService(userService UserService) SecurityService {
	return &securityService{
		userService: userService,
	}
}

func (service *securityService) IsAdmin(chatID int64) (bool, error) {
	user, err := service.userService.FindUserByChatID(chatID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}

		return false, err
	}

	return user.Role == model.UserRoleAdmin, nil
}
