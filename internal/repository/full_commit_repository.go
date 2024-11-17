package repository

import (
	"database/sql"

	"github.com/epistax1s/gomer/internal/database"
	"github.com/epistax1s/gomer/internal/model"
	"gorm.io/gorm"
)

type FullCommitRepository interface {
	FindAllByDate(*database.Date) ([]model.FullCommit, error)
}

type fullCommitRepository struct {
	db *gorm.DB
}

func NewFullCommit(db *gorm.DB) FullCommitRepository {
	return &fullCommitRepository{
		db: db,
	}
}

func (repo *fullCommitRepository) FindAllByDate(date *database.Date) ([]model.FullCommit, error) {
	var fullCommits []model.FullCommit
	result := repo.db.Raw(`
		SELECT 
			us.name, 
			us.username,
			us.chat_id, 
			dp.id as department_id, 
			dp.department_name, 
			CASE 
				WHEN cm.id IS NULL THEN FALSE
				ELSE TRUE
			END as commit_sent,
			cm.commit_payload
		FROM "tg_user" us
		JOIN "department" dp ON us.department_id = dp.id
		LEFT JOIN "commit" cm ON us.id = cm.user_id AND cm.commit_date = @commit_date
		WHERE us.status = @user_status
	`, sql.Named("commit_date", date), sql.Named("user_status", model.UserStatusActive)).Scan(&fullCommits)

	return fullCommits, result.Error
}
