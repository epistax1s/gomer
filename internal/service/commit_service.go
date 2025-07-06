package service

import (
	"errors"
	"github.com/epistax1s/gomer/internal/database"
	"github.com/epistax1s/gomer/internal/log"
	"github.com/epistax1s/gomer/internal/model"
	"github.com/epistax1s/gomer/internal/repository"
	"gorm.io/gorm"
)

type CommitService interface {
	FindByUserIDAndDate(userID int64, date *database.Date) (*model.Commit, error)
	FindByChatIDAndDate(chatID int64, date *database.Date) (*model.Commit, error)
	FindAllByDate(date *database.Date) ([]model.Commit, error)
	CreateCommit(chatID int64, payload string, data *database.Date) error
	UpdateCommit(id int64, payload string) (*model.Commit, error)
}

type commitService struct {
	userServices UserService
	commitRepo   repository.CommitRepository
}

func NewCommitService(
	userService UserService,
	commitRepo repository.CommitRepository) CommitService {

	return &commitService{
		userServices: userService,
		commitRepo:   commitRepo,
	}
}

func (service *commitService) FindByUserIDAndDate(userID int64, date *database.Date) (*model.Commit, error) {
	commit, err := service.commitRepo.FindByUserIDAndDate(userID, date)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Debug(
				"commit not found",
				"userID", userID, "commitDate", date)

			return nil, nil
		}

		log.Error(
			"commit search, error when searching for a commit by userID and date",
			"userID", userID, "commitDate", date)

		return nil, err
	}

	return commit, nil
}

func (service *commitService) FindByChatIDAndDate(chatID int64, date *database.Date) (*model.Commit, error) {
	user, err := service.userServices.FindUserByChatID(chatID)
	if err != nil {
		log.Error(
			"commit search, error when searching for a user by chatID",
			"chatID", chatID, "err", err)

		return nil, err
	}

	if user == nil {
		return nil, nil
	}

	return service.FindByUserIDAndDate(user.ID, date)
}

func (service *commitService) FindAllByDate(date *database.Date) ([]model.Commit, error) {
	return service.commitRepo.FindAllByDate(date)
}

func (service *commitService) CreateCommit(chatID int64, payload string, date *database.Date) error {
	user, err := service.userServices.FindUserByChatID(chatID)
	if err != nil {
		log.Error(
			"creating a commit, error when searching for a user",
			"chatID", chatID, "err", err)

		return errors.New("error when creating a commit")
	}

	if user == nil {
		log.Error(
			"creating a commit, the user was not found",
			"chatID", chatID)

		return errors.New("error when creating a commit")
	}

	err = service.commitRepo.Create(&model.Commit{
		User:    *user,
		Payload: payload,
		Date:    date,
	})

	if err != nil {
		log.Error(
			"error when creating a commit in the database",
			"chatID", chatID, "err", err)

		return errors.New("error when creating a commit")
	}

	log.Info(
		"the commit was successfully created",
		"chatID", chatID, "date", date)

	return nil
}

func (service *commitService) UpdateCommit(id int64, payload string) (*model.Commit, error) {
	commit, err := service.commitRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error(
				"attempt to update the commit, the commit was not found",
				"id", id)

			return nil, err
		}

		log.Error(
			"attempt to update the commit, error when trying to get a commit from the database",
			"id", id, "err", err)

		return nil, err
	}

	commit.Payload = payload
	if err := service.commitRepo.Update(commit); err != nil {
		log.Error(
			"error when updating a commit in the database",
			"id", id, "commit", commit, "err", err)

		return nil, err
	}

	log.Info(
		"the commit has been successfully updated",
		"id", id)

	return commit, nil
}
