package model

type FullCommit struct {
	Name           string `gorm:"column:name"`
	Username       string `gorm:"column:username"`
	ChatID         int64  `gorm:"column:chat_id"`
	DepartmentID   int64  `gorm:"column:department_id"`
	DepartmentName string `gorm:"column:department_name"`
	CommitSent     bool   `gorm:"column:commit_sent"`
	CommitPayload  string `gorm:"column:commit_payload"`
}

type ByDepartamentAndName []FullCommit

func (a ByDepartamentAndName) Len() int {
	return len(a)
}

func (a ByDepartamentAndName) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByDepartamentAndName) Less(i, j int) bool {
	if a[i].DepartmentID == a[j].DepartmentID {
		return a[i].Name < a[j].Name
	}

	return a[i].DepartmentID < a[j].DepartmentID
}
