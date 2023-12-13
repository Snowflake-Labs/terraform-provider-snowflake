package sdk

import (
	"time"
)

// fix timestamp merge
func ParseTimestampWithOffset(s string, dateTimeFormat string) (string, error) {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return err.Error(), err
	}
	_, offset := t.Zone()
	adjustedTime := t.Add(-time.Duration(offset) * time.Second)
	adjustedTimeFormat := adjustedTime.Format(dateTimeFormat)
	return adjustedTimeFormat, nil
}
