package acceptance

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func IsGreaterOrEqualTo(greaterOrEqualValue int) resource.CheckResourceAttrWithFunc {
	return func(value string) error {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("unable to parse value %s as integer, err = %w", value, err)
		}

		if intValue < greaterOrEqualValue {
			return fmt.Errorf("expected value %d to be greater or equal to %d", intValue, greaterOrEqualValue)
		}

		return nil
	}
}
