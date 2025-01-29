package utils

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/epistax1s/gomer/internal/log"
)

type Callback interface {
	GetType() string
}
type BaseCallback struct {
	Type string `json:"type"`
}

func (b *BaseCallback) GetType() string {
	return b.Type
}

/*
Common callbacks.
*/
const (
	EXIT = "exit"
)

type ExitCallback struct {
	BaseCallback
}

func NewExitCallback() string {
	return Encode(
		&ExitCallback{
			BaseCallback: BaseCallback{
				Type: EXIT,
			},
		},
	)
}

/*
Group administration callbacks (Admin Group management).
*/
const (
	AG_PREV       = "ag_prev"
	AG_NEXT       = "ag_next"
	AG_SELECT     = "ag_select"
	AG_UNLINK     = "ag_unlink"
	AG_GROUP_LIST = "ag_group_list"
)

type AGPrevCallback struct {
	BaseCallback
	Page int
}

func NewAGPrevCallback(page int) string {
	return Encode(
		&AGPrevCallback{
			BaseCallback: BaseCallback{
				Type: AG_PREV,
			},
			Page: page,
		},
	)
}

type AGNextCallback struct {
	BaseCallback
	Page int
}

func NewAGNextCallback(page int) string {
	return Encode(
		&AGNextCallback{
			BaseCallback: BaseCallback{
				Type: AG_NEXT,
			},
			Page: page,
		},
	)
}

type AGSelectCallback struct {
	BaseCallback
	Page    int
	GroupID int
}

func NewAGSelectCallback(page int, groupID int) string {
	return Encode(
		&AGSelectCallback{
			BaseCallback: BaseCallback{
				Type: AG_SELECT,
			},
			Page:    page,
			GroupID: groupID,
		},
	)
}

type AGUnlinkCallback struct {
	BaseCallback
	GroupID int
}

func NewAGUnlinkCallback(groupID int) string {
	return Encode(
		&AGUnlinkCallback{
			BaseCallback: BaseCallback{
				Type: AG_UNLINK,
			},
			GroupID: groupID,
		},
	)
}

type AGGroupListCallback struct {
	BaseCallback
	Page int
}

func NewAGGroupListCallback(page int) string {
	return Encode(
		&AGGroupListCallback{
			BaseCallback: BaseCallback{
				Type: AG_GROUP_LIST,
			},
			Page: page,
		},
	)
}

/*
Calendar callbacks.
*/

const (
	CALENDAR_PREV = "cl_prev"
	CALENDAR_NEXT = "cl_next"
	CALENDAR_DATE = "cl_date"
)

type CalendarPrevCallback struct {
	BaseCallback
	Year  int
	Month time.Month
}

func NewCalendarPrevCallback(year int, month time.Month) string {
	return Encode(
		&CalendarPrevCallback{
			BaseCallback: BaseCallback{
				Type: CALENDAR_PREV,
			},
			Year:  year,
			Month: month,
		},
	)
}

type CalendarNextCallback struct {
	BaseCallback
	Year  int
	Month time.Month
}

func NewCalendarNextCallback(year int, month time.Month) string {
	return Encode(
		&CalendarNextCallback{
			BaseCallback: BaseCallback{
				Type: CALENDAR_NEXT,
			},
			Year:  year,
			Month: month,
		},
	)
}

type CalendarDateCallback struct {
	BaseCallback
	Date string
}

func NewCalendarDateCallback(date string) string {
	return Encode(
		&CalendarDateCallback{
			BaseCallback: BaseCallback{
				Type: CALENDAR_DATE,
			},
			Date: date,
		},
	)
}

func Encode(с Callback) string {
	data, err := json.Marshal(с)
	if err != nil {
		log.Error(err.Error())
		return ""
	}
	return string(data)
}

func Decode(data string) (Callback, error) {
	var base BaseCallback
	if err := json.Unmarshal([]byte(data), &base); err != nil {
		return nil, err
	}

	// Replace it with a reflection, and in the meantime, let's pretend you didn't see it.
	switch base.Type {
	case EXIT:
		var result ExitCallback
		err := json.Unmarshal([]byte(data), &result)
		return &result, err
	case AG_PREV:
		var result AGPrevCallback
		err := json.Unmarshal([]byte(data), &result)
		return &result, err
	case AG_NEXT:
		var result AGNextCallback
		err := json.Unmarshal([]byte(data), &result)
		return &result, err
	case AG_SELECT:
		var result AGSelectCallback
		err := json.Unmarshal([]byte(data), &result)
		return &result, err
	case AG_GROUP_LIST:
		var result AGGroupListCallback
		err := json.Unmarshal([]byte(data), &result)
		return &result, err
	case AG_UNLINK:
		var result AGUnlinkCallback
		err := json.Unmarshal([]byte(data), &result)
		return &result, err
	case CALENDAR_PREV:
		var result CalendarPrevCallback
		err := json.Unmarshal([]byte(data), &result)
		return &result, err
	case CALENDAR_NEXT:
		var result CalendarNextCallback
		err := json.Unmarshal([]byte(data), &result)
		return &result, err
	case CALENDAR_DATE:
		var result CalendarDateCallback
		err := json.Unmarshal([]byte(data), &result)
		return &result, err
	}

	return nil, fmt.Errorf("unsupported callback type: %s", base.Type)
}
