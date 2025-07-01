package redmine

type TimeEntriesResponse struct {
	TimeEntries []TimeEntry `json:"time_entries"`
	TotalCount  int         `json:"total_count"`
	Offset      int         `json:"offset"`
	Limit       int         `json:"limit"`
}

type TimeEntry struct {
	ID       int     `json:"id"`
	Project  Project `json:"project"`
	Issue    Issue   `json:"issue"`
	User     User    `json:"user"`
	Activity Named   `json:"activity"`
	Hours    float64 `json:"hours"`
	Comments string  `json:"comments"`
	SpentOn  string  `json:"spent_on"`
}

type Project struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Issue struct {
	ID int `json:"id"`
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Named struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
