package model

// TODO пробежаться по всем сущностям и посмотреть что там с foreignkey связями
type User struct {
	ID           int64  `gorm:"column:id;primaryKey"`
	DepartmentId int64  `gorm:"column:department_id;foreignkey:user_department_fk;references:id"`
	Order        int64  `gorm:"column:order"`
	ChatID       int64  `gorm:"column:chat_id"`
	Name         string `gorm:"column:name"`
	Username     string `gorm:"column:username"`
	Role         string `gorm:"column:role"`
	Status       string `gorm:"column:status"`
}

const (
	UserTable          = "tg_user"
	UserChatIDColumn   = "chat_id"
	UserUsernameColumn = "username"
)

const (
	UserStatusActive  = "active"
	UserStatusDeleted = "deleted"
)

const (
	UserRoleUser  = "USER"
	UserRoleAdmin = "ADMIN"
)

func (User) TableName() string {
	return UserTable
}


