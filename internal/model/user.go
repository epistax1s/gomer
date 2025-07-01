package model

// TODO пробежаться по всем сущностям и посмотреть что там с foreignkey связями
type User struct {
	ID           int64  `gorm:"column:id;primaryKey"`
	ChatID       int64  `gorm:"column:chat_id"`
	Name         string `gorm:"column:name"`
	Username     string `gorm:"column:username"`
	DepartmentId int64  `gorm:"column:department_id;foreignkey:user_department_fk;references:id"`
	Order        int64  `gorm:"column:order"`
	Role         string `gorm:"column:role"`
	Status       string `gorm:"column:status"`
	CommitSrc    string `gorm:"column:commit_src"`
}

const (
	UserTable          = "tg_user"
	UserChatIDColumn   = "chat_id"
	UserUsernameColumn = "username"
)

const (
	UserStatusLimbo   = "limbo"
	UserStatusActive  = "active"
	UserStatusDeleted = "deleted"
	UserStatusSystem  = "system"
)

const (
	UserRoleUser  = "USER"
	UserRoleAdmin = "ADMIN"
)

const (
	UserCommitSrcManual     = "manual"
	UserCommitSrcRedmine    = "redmine"
	UserCommitSrcRedmineExt = "redmin_ext"
)

func (User) TableName() string {
	return UserTable
}
