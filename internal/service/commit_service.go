package service

import (
	"errors"

	"gorm.io/gorm"

	"github.com/epistax1s/gomer/internal/database"
	"github.com/epistax1s/gomer/internal/log"
	"github.com/epistax1s/gomer/internal/model"
	"github.com/epistax1s/gomer/internal/repository"
)

type CommitService interface {
	FindCommitByUserIdAndDate(int64, *database.Date) (*model.Commit, error)
	FindAllCommitsByDate(*database.Date) ([]model.Commit, error)
	CreateCommit(int64, string, *database.Date) error
	UpdateCommit(int64, string) (*model.Commit, error)
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

func (service *commitService) FindCommitByUserIdAndDate(chatID int64, date *database.Date) (*model.Commit, error) {
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

	commit, err := service.commitRepo.FindByUserIdAndDate(user.ID, date)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Info(
				"commit not found",
				"userID", user.ID, "commitDate", date)

			return nil, nil
		}

		log.Error(
			"commit search, error when searching for a commit by userID and date",
			"userID", user.ID, "commitDate", date)

		return nil, err
	}

	return commit, nil
}

func (service *commitService) FindAllCommitsByDate(date *database.Date) ([]model.Commit, error) {
	return service.commitRepo.FindAllByDate(date)
}

func (service *commitService) CreateCommit(chatID int64, payload string, commitDate *database.Date) error {
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
		Date:    commitDate,
	})

	if err != nil {
		log.Error(
			"error when creating a commit in the database",
			"chatID", chatID, "err", err)

		return errors.New("error when creating a commit")
	}

	log.Info(
		"the commit was successfully created",
		"chatID", chatID, "commitDate", commitDate)

	return nil
}

func (service *commitService) UpdateCommit(commitID int64, payload string) (*model.Commit, error) {
	commit, err := service.commitRepo.FindById(commitID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error(
				"attempt to update the commit, the commit was not found",
				"commitID", commit)

			return nil, err
		}

		log.Error(
			"attempt to update the commit, error when trying to get a commit from the database",
			"commitID", commit, "err", err)

		return nil, err
	}

	commit.Payload = payload
	if err := service.commitRepo.Update(commit); err != nil {
		log.Error(
			"error when updating a commit in the database",
			"commitID", commitID, "commit", commit, "err", err)

		return nil, err
	}

	log.Info(
		"the commit has been successfully updated",
		"commitID", commitID)

	return commit, nil
}
