package callback

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
	PREV           = "prev"    // pagination
	NEXT           = "next"    // pagination
	SWAP           = "swap"    // swap 2 entities
	ADD            = "add"     // add new entity
	SELECT         = "select"  // select entity
	DELETE         = "delete"  // delete entity
	BACK_TO_LIST   = "bt_list" // back to the list of entities
	EXIT           = "exit"
	ACTION_CONFIRM = "a_confirm"
	ACTION_CANCEL  = "a_cancel"
)

type PrevCallback struct {
	BaseCallback
	Page int
}

func NewPrevCallback(page int) string {
	return Encode(
		&PrevCallback{
			BaseCallback: BaseCallback{
				Type: PREV,
			},
			Page: page,
		},
	)
}

type NextCallback struct {
	BaseCallback
	Page int
}

func NewNextCallback(page int) string {
	return Encode(
		&NextCallback{
			BaseCallback: BaseCallback{
				Type: NEXT,
			},
			Page: page,
		},
	)
}

type SwapCallback struct {
	BaseCallback
	EntityIDA int64 `json:"a"`
	EntityIDB int64 `json:"b"`
}

func NewSwapCallback(entityAID int64, entityBID int64) string {
	return Encode(
		&SwapCallback{
			BaseCallback: BaseCallback{
				Type: NEXT,
			},
			EntityIDA: entityAID,
			EntityIDB: entityBID,
		},
	)
}

type SelectCallback struct {
	BaseCallback
	EntityID int
	Page     int
}

func NewSelectCallback(entityID int, page int) string {
	return Encode(
		&SelectCallback{
			BaseCallback: BaseCallback{
				Type: NEXT,
			},
			EntityID: entityID,
			Page:     page,
		},
	)
}

type DeleteCallback struct {
	BaseCallback
	EntityID int
}

func NewDeleteCallback(entityID int) string {
	return Encode(
		&DeleteCallback{
			BaseCallback: BaseCallback{
				Type: DELETE,
			},
			EntityID: entityID,
		},
	)
}

type AddCallback struct {
	BaseCallback
}

func NewAddCallback() string {
	return Encode(
		&AddCallback{
			BaseCallback: BaseCallback{
				Type: ADD,
			},
		},
	)
}

type ListCallback struct {
	BaseCallback
	Page int
}

func NewListCallback(page int) string {
	return Encode(
		&ListCallback{
			BaseCallback: BaseCallback{
				Type: BACK_TO_LIST,
			},
			Page: page,
		},
	)
}

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
Users administration callbacks (Admin Group management).
*/
const (
	AU_PREV      = "au_prev"
	AU_NEXT      = "au_next"
	AU_MOVE_UP   = "au_mvup"
	AU_MOVE_DOWN = "au_mvdw"
	AU_SELECT    = "au_select"
	AU_USER_LIST = "au_user_list"
)

type AUPrevCallback struct {
	BaseCallback
	Page int
}

/*
Groups administration callbacks (Admin Group management).
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

func NewActionConfirmCallback() string {
	return Encode(
		&BaseCallback{
			Type: ACTION_CONFIRM,
		},
	)
}

func NewActionCancelCallback() string {
	return Encode(
		&BaseCallback{
			Type: ACTION_CANCEL,
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
	case PREV:
		var result PrevCallback
		err := json.Unmarshal([]byte(data), &result)
		return &result, err
	case NEXT:
		var result NextCallback
		err := json.Unmarshal([]byte(data), &result)
		return &result, err
	case SWAP:
		var result SwapCallback
		err := json.Unmarshal([]byte(data), &result)
		return &result, err
	case ADD:
		var result AddCallback
		err := json.Unmarshal([]byte(data), &result)
		return &result, err
	case SELECT:
		var result SelectCallback
		err := json.Unmarshal([]byte(data), &result)
		return &result, err
	case DELETE:
		var result DeleteCallback
		err := json.Unmarshal([]byte(data), &result)
		return &result, err
	case BACK_TO_LIST:
		var result ListCallback
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
	case ACTION_CONFIRM:
		var result BaseCallback
		err := json.Unmarshal([]byte(data), &result)
		return &result, err
	case ACTION_CANCEL:
		var result BaseCallback
		err := json.Unmarshal([]byte(data), &result)
		return &result, err
	}

	return nil, fmt.Errorf("unsupported callback type: %s", base.Type)
}
