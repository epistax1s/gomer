package model

// TODO пробежаться по всем сущностям и посмотреть что там с foreignkey связями
type User struct {
	ID           int64      `gorm:"column:id;primaryKey"`
	ChatID       int64      `gorm:"column:chat_id"`
	RedmineID    *int64      `gorm:"column:redmine_id"`
	Name         string     `gorm:"column:name"`
	Username     string     `gorm:"column:username"`
	DepartmentID int64      `gorm:"column:department_id;foreignkey:user_department_fk;references:id"`
	Department   Department `gorm:"foreignKey:DepartmentID"`
	Order        int64      `gorm:"column:order"`
	Role         string     `gorm:"column:role"`
	Status       string     `gorm:"column:status"`
	CommitSrc    string     `gorm:"column:commit_src"`
	EE           bool      	`gorm:"column:ee"`
}

const (
	UserTable          = "tg_user"
	UserChatIDColumn   = "chat_id"
	UserUsernameColumn = "username"
	UserStatusColumn   = "status"
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

type ByDepartmentOrderAndName []User

func (a ByDepartmentOrderAndName) Len() int {
	return len(a)
}

func (a ByDepartmentOrderAndName) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByDepartmentOrderAndName) Less(i, j int) bool {
	// If Department.Order is equal, sorted by Name
	if a[i].Department.Order == a[j].Department.Order {
		return a[i].Name < a[j].Name
	}

	// Otherwise, we sort by Department.Order
	return a[i].Department.Order < a[j].Department.Order
}
