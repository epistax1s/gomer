package model

import "github.com/epistax1s/gomer/internal/database"

type Commit struct {
	ID      int64          `gorm:"column:id;primaryKey"`
	UserID  int64          `gorm:"column:user_id"`
	User    User           `gorm:"foreignkey:UserID;references:ID"`
	Payload string         `gorm:"column:commit_payload"`
	Date    *database.Date `gorm:"column:commit_date"`
}

const (
	CommitTable        = "commit"
	CommitUserIDColumn = "user_id"
	CommitDateColumn   = "commit_date"
)

func (Commit) TableName() string {
	return CommitTable
}
