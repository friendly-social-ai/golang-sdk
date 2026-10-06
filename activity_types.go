package sdk

import "encoding/json"

// --- ACTIVITY ID ---

// ActivityId represents the unique identifier of activity.
type ActivityId struct {
	value int64
}

// NewActivityId creates new ActivityId from int64.
func NewActivityId(i int64) ActivityId {
	return ActivityId{value: i}
}

// Value returns ActivityId as a plain int64.
func (i ActivityId) Value() int64 {
	return i.value
}

func (i ActivityId) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.value)
}

func (i *ActivityId) UnmarshalJSON(bytes []byte) error {
	return json.Unmarshal(bytes, &i.value)
}
