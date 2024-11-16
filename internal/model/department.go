package model

type Department struct {
	ID    int64  `gorm:"column:id;primaryKey"`
	Order int64  `gorm:"column:order"`
	Name  string `gorm:"column:department_name"`
}

func (Department) TableName() string {
	return "department"
}
