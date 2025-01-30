package repository

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/epistax1s/gomer/internal/model"
)

type UserRepository interface {
	Create(user *model.User) error
	Update(user *model.User) error
	Delete(id int64) (bool, error)
	FindPaginated(page int, pageSize int) ([]model.User, error)
	FindAll() ([]model.User, error)
	FindByUsername(username string) (*model.User, error)
	FindByChatID(chatID int64) (*model.User, error)
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

func (repo *userRepository) Delete(id int64) (bool, error) {
	result := repo.db.Delete(&model.User{}, id)
	return result.RowsAffected > 0, result.Error
}

func (repo *userRepository) FindPaginated(page int, pageSize int) ([]model.User, error) {
	var users []model.User

	offset := (page - 1) * pageSize

	result := repo.db.Offset(offset).Limit(pageSize).Find(&users)
	return users, result.Error
}

func (repo *userRepository) FindAll() ([]model.User, error) {
	var users []model.User
	result := repo.db.Find(&users)
	return users, result.Error
}

func (repo *userRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User

	result := repo.db.
		Where(fmt.Sprintf("%s = ?", model.UserUsernameColumn), username).
		First(&user)

	return &user, result.Error
}

func (repo *userRepository) FindByChatID(chatID int64) (*model.User, error) {
	var user model.User

	result := repo.db.
		Where(fmt.Sprintf("%s = ?", model.UserChatIDColumn), chatID).
		First(&user)

	return &user, result.Error
}

func (repo *userRepository) CountAll() (int64, error) {
	var count int64
	result := repo.db.Model(&model.User{}).Count(&count)
	return count, result.Error
}
