package sdk

import (
	"fmt"
)

func IsValidDataType(v string) bool {
	dt := DataTypeFromString(v)
	return dt != DataTypeUnknown
}

func checkExclusivePointers(ptrs []interface{}) error {
	count := 0
	for _, v := range ptrs {
		// Types differ so we can't directly compare to `nil`
		if fmt.Sprintf("%v", v) != "<nil>" {
			count++
		}
	}
	if count != 1 {
		return fmt.Errorf("%d values set", count)
	}
	return nil
}
