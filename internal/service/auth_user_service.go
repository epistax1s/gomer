package service

import (
	"github.com/epistax1s/gomer/internal/model"
	"github.com/epistax1s/gomer/internal/repository"
)

type AuthUserService interface {
	Register(int64, string) error
	IsRegistered(int64) bool
}

type authUserService struct {
	authUserRepo repository.AuthUserRepository
}

func NewAuthUserService(authUserRepo repository.AuthUserRepository) AuthUserService {
	return &authUserService{
		authUserRepo: authUserRepo,
	}
}

func (service *authUserService) Register(chatID int64, username string) error {
	authUser := &model.AuthUser{
		ChatID:   chatID,
		Username: username,
	}

	return service.authUserRepo.Create(authUser)
}

func (service *authUserService) IsRegistered(chatID int64) bool {
	user, err := service.authUserRepo.FindByChatID(chatID)
	if err != nil || user == nil {
		return false
	} else {
		return true
	}
}
