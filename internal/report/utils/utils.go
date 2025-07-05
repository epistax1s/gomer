package utils

import (
	"regexp"
	"strings"
	"time"

	"github.com/epistax1s/gomer/internal/database"
	"github.com/epistax1s/gomer/internal/log"
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

func ShouldExcludeComment(comment string, excludePatterns [] string) bool {
	if comment == "" {
		return true
	}

	if len(excludePatterns) == 0 {
		return false
	}

	for _, pattern := range excludePatterns {
		if pattern == "" {
			continue
		}

		matched, err := regexp.MatchString(pattern, strings.ToLower(comment))
		if err != nil {
			log.Error("Invalid regex pattern", "pattern", pattern, "error", err)
			continue
		}

		if matched {
			return true
		}
	}

	return false
}