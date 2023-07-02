package schema

import (
	"encoding/json"
	"fmt"
	"time"
)

type CreatedTime time.Time

func (t CreatedTime) MarshalJSON() ([]byte, error) {
	value := time.Time(t)
	//created, err := time.Parse(time.RFC3339, d.created_at.String)

	bytes, err := json.Marshal(value.Format(time.RFC3339))
	if err != nil {
		return nil, fmt.Errorf("error marshal  CreatedTime %v", value)
	}
	return bytes, nil
}
func (t *CreatedTime) UnmarshalJSON(b []byte) error {
	var createdTimeString string
	err := json.Unmarshal(b, &createdTimeString)
	if err != nil {
		return fmt.Errorf("error unmarshal  CreatedTime %v", b)
	}
	createdTime, err := time.Parse(time.RFC3339, createdTimeString)
	if err != nil {
		return fmt.Errorf("error parse  CreatedTime to RFC3339 %v", b)
	}
	*t = CreatedTime(createdTime)
	return nil
}
