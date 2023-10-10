package sdk_integration_tests

import "time"

func ParseTimestampWithOffset(s string) (*time.Time, error) {
	t, err := time.Parse("2006-01-02T15:04:05-07:00", s)
	if err != nil {
		return nil, err
	}
	_, offset := t.Zone()
	adjustedTime := t.Add(-time.Duration(offset) * time.Second)
	return &adjustedTime, nil
}
