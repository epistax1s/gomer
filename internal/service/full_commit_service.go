package service

import (
	"github.com/epistax1s/gomer/internal/database"
	"github.com/epistax1s/gomer/internal/log"
	"github.com/epistax1s/gomer/internal/model"
	"github.com/epistax1s/gomer/internal/repository"
)

type FullCommitService interface {
	FindAllByDate(*database.Date) ([]model.FullCommit, error)
}

type fullCommitService struct {
	fullCommitRepo repository.FullCommitRepository
}

func NewFullCommitService(fullCommitRepo repository.FullCommitRepository) FullCommitService {
	return &fullCommitService{
		fullCommitRepo: fullCommitRepo,
	}
}

func (service fullCommitService) FindAllByDate(date *database.Date) ([]model.FullCommit, error) {
	fullCommits, err := service.fullCommitRepo.FindAllByDate(date)
	if err != nil {
		log.Error(
			"error when searching for full commits for the specified date",
			"date", date, "err", err)
		return nil, err
	}
	return fullCommits, nil
}
