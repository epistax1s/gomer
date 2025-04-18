package service

import (
	"github.com/epistax1s/gomer/internal/repository"
)

type AuthKeyService interface {
	IsValidKey(string) bool
}

type authKeyService struct {
	authKeyRepo repository.AuthKeyRepository
}

func NewAuthKeyService(authKeyRepo repository.AuthKeyRepository) AuthKeyService {
	return &authKeyService{
		authKeyRepo: authKeyRepo,
	}
}

func (service *authKeyService) IsValidKey(key string) bool {
	_, err := service.authKeyRepo.FindByKey(key)
	if err == nil {
		return true
	} else {
		return false
	}
}
