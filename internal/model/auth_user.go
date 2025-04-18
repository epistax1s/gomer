package model

type AuthUser struct {
	ID       int64  `gorm:"column:id;primaryKey"`
	ChatID   int64  `gorm:"column:chat_id"`
	Username string `gorm:"column:username"`
}

const (
	AuthUserTable          = "auth_user"
	AuthUserChatIDColumn   = "chat_id"
	AuthUserUsernameColumn = "username"
)

func (AuthUser) TableName() string {
	return AuthUserTable
}
