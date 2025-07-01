package model

type Invitation struct {
	ID          int64   `gorm:"column:id;primaryKey;autoIncrement"`
	Code        string  `gorm:"column:code;unique;not null"`
	CreatedByID int64   `gorm:"column:created_by_user_id;not null"`
	CreatedBy   User    `gorm:"foreignKey:CreatedByID;references:ID"`
	CreatedAt   string  `gorm:"column:created_at;not null"`
	Used        bool    `gorm:"column:used;not null;default:true"`
	UsedByID    *int64  `gorm:"column:used_by_user_id"`
	UsedBy      *User   `gorm:"foreignKey:UsedByID;references:ID"`
	UsedAt      *string `gorm:"column:used_at"`
}

const (
	InvitationTable = "invitation"
)

func (Invitation) TableName() string {
	return InvitationTable
}
