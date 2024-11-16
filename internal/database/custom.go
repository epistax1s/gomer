package database

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type Date struct {
	time.Time
}

func (date *Date) Scan(value interface{}) error {
	if v, ok := value.(string); ok {
		t, err := time.Parse("2006-01-02", v)
		if err != nil {
			return err
		}
		date.Time = t
		return nil
	}
	return fmt.Errorf("incorrect data type: %T", value)
}

func (date Date) Value() (driver.Value, error) {
	return date.Time.Format("2006-01-02"), nil
}

func (date Date) String() string {
	return date.Time.Format("2006-01-02")
}
