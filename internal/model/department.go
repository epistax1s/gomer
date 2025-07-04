package model

type Department struct {
	ID    int64  `gorm:"column:id;primaryKey"`
	Order int64  `gorm:"column:order"`
	Name  string `gorm:"column:department_name"`
	Type  string `gorm:"column:type"`
}

const (
	DepartmentTable      = "department"
	DepartmentTypeColumn = "type"
)

const (
	DepartmentTypeSystem   = "system"
	DepartmentTypeNormal   = "normal"
)

func (Department) TableName() string {
	return DepartmentTable
}
