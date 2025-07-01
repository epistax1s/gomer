package utils

import (
	"time"

	"github.com/epistax1s/gomer/internal/database"
)

func IsBusinessDay(t *time.Time) bool {
	return t.Weekday() != time.Saturday && t.Weekday() != time.Sunday
}

func CastToBuildDate(t *time.Time) *database.Date {
	if t.Weekday() == time.Monday {
		return &database.Date{Time: t.Add(-72 * time.Hour)}
	}
	return &database.Date{Time: t.Add(-24 * time.Hour)}
}
