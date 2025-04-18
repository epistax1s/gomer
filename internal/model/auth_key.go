package model

type AuthKey struct {
	ID  int64  `gorm:"column:id;primaryKey"`
	Key string `gorm:"column:key"`
}

const (
	AuthKeyTable  = "auth_key"
	AuthKeyColumn = "key"
)

func (AuthKey) TableName() string {
	return AuthKeyTable
}
