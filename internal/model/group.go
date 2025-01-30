package model

type Group struct {
	ID      int64  `gorm:"column:id;primaryKey"`
	GroupID int64  `gorm:"column:group_id"`
	Title   string `gorm:"column:title"`
}

func (Group) TableName() string {
	return "tg_group"
}

const (
	GroupGroupID = "group_id"
)
