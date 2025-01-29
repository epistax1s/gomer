package utils

import "encoding/json"

type Callback interface {
	GetType() string
}
type BaseCallback struct {
	Type string `json:"type"`
}

func (b *BaseCallback) GetType() string {
	return b.Type
}

const (
	AG_PREV   = "ag_pv"
	AG_NEXT   = "ag_nt"
	AG_SELECT = "ag_sl"
	AG_MENU   = "ag_mn"
	AG_UNLINK = "ag_ul"
)

type AGPrevCallback struct {
	BaseCallback
}

type AGNextCallback struct {
	BaseCallback
}
type AGSelectCallback struct {
	BaseCallback
}
type AGMenuCallback struct {
	BaseCallback
}
type AGUnlinkCallback struct {
	BaseCallback
}

func ParseCallback(data string) (Callback, error) {
	var base BaseCallback
	if err := json.Unmarshal([]byte(data), &base); err != nil {
		return nil, err
	}

	switch base.Type {
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
	case AG_MENU:
		var result AGMenuCallback
		err := json.Unmarshal([]byte(data), &result)
		return &result, err
	case AG_UNLINK:
		var result AGUnlinkCallback
		err := json.Unmarshal([]byte(data), &result)
		return &result, err
	}

	return nil, nil
}
