package acceptance

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsGreaterOrEqualTo(t *testing.T) {
	testCases := []struct {
		Name             string
		GreaterOrEqualTo int
		Actual           string
		Error            string
	}{
		{
			Name:             "validation: smaller than expected",
			GreaterOrEqualTo: 20,
			Actual:           "10",
			Error:            "expected value greater or equal to 20, got 10",
		},
		{
			Name:             "validation: not int value",
			GreaterOrEqualTo: 20,
			Actual:           "not_int",
			Error:            "unable to parse value not_int as integer, err = strconv.Atoi: parsing \"not_int\": invalid syntax",
		},
		{
			Name:             "validation: equal value",
			GreaterOrEqualTo: 20,
			Actual:           "20",
		},
		{
			Name:             "validation: greater value",
			GreaterOrEqualTo: 20,
			Actual:           "30",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			err := IsGreaterOrEqualTo(testCase.GreaterOrEqualTo)(testCase.Actual)
			if testCase.Error != "" {
				assert.ErrorContains(t, err, testCase.Error)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
