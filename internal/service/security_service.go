package service

import (
	"github.com/epistax1s/gomer/internal/model"
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
		return false, err
	}

	if user == nil {
		return false, nil
	}

	return user.Role == model.UserRoleAdmin, nil
}
