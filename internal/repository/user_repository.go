package repository

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/epistax1s/gomer/internal/model"
)

type UserRepository interface {
	Create(user *model.User) error
	Update(user *model.User) error
	FindByID(id int64) (*model.User, error)
	FindByChatID(chatID int64) (*model.User, error)
	FindByUsername(username string) (*model.User, error)
	FindAll() ([]model.User, error)
	FindAllActive() ([]model.User, error)
	CountAll() (int64, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (repo *userRepository) Create(user *model.User) error {
	return repo.db.Create(user).Error
}

func (repo *userRepository) Update(user *model.User) error {
	return repo.db.Updates(user).Error
}

func (repo *userRepository) FindByID(id int64) (*model.User, error) {
	var user model.User
	result := repo.db.First(&user, id)
	return &user, result.Error
}

func (repo *userRepository) FindByChatID(chatID int64) (*model.User, error) {
	var user model.User

	result := repo.db.
		Where(fmt.Sprintf("%s = ?", model.UserChatIDColumn), chatID).
		First(&user)

	return &user, result.Error
}

func (repo *userRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User

	result := repo.db.
		Where(fmt.Sprintf("%s = ?", model.UserUsernameColumn), username).
		First(&user)

	return &user, result.Error
}

func (repo *userRepository) FindAll() ([]model.User, error) {
	var users []model.User
	
	result := repo.db.
		Preload("Department").
		Find(&users)

	return users, result.Error
}

func (repo *userRepository) FindAllActive() ([]model.User, error) {
	var users []model.User

	result := repo.db.
		Preload("Department").
		Where(fmt.Sprintf("%s = ?", model.UserStatusColumn), model.UserStatusActive).
		Find(&users)

	return users, result.Error
}

func (repo *userRepository) CountAll() (int64, error) {
	var count int64
	result := repo.db.Model(&model.User{}).Count(&count)
	return count, result.Error
}
