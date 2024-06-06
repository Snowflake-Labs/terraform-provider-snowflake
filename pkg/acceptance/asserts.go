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
			return err
		}

		if intValue < greaterOrEqualValue {
			return fmt.Errorf("expected value greater or equal to %d, got %d", greaterOrEqualValue, intValue)
		}

		return nil
	}
}
