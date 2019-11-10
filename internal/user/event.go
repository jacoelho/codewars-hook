package user

import (
	"encoding/json"
	"fmt"
)

// Event represents webhook event user message
type Event struct {
	HonorUpgraded *HonorUpgraded
}

// HonorUpgraded represents an honor upgraded event
type HonorUpgraded struct {
	ID         string `json:"id"`
	Honor      int    `json:"honor"`
	HonorDelta int    `json:"honor_delta"`
}

// UnmarshalJSON implements json.Unmarshaler interface
func (event *Event) UnmarshalJSON(data []byte) error {
	var payload struct {
		Action string           `json:"action"`
		Event  *json.RawMessage `json:"user"`
	}

	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}

	switch payload.Action {
	case "honor_changed":
		var data HonorUpgraded

		err := json.Unmarshal(*payload.Event, &data)

		event.HonorUpgraded = &data
		return err
	default:
		return fmt.Errorf("unknown action %s", payload.Action)
	}
}
