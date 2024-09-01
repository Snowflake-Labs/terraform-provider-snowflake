package resources_test

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// checkBool is deprecated and will be removed with resources rework (and replaced with new assertions)
func checkBool(path, attr string, value bool) func(*terraform.State) error {
	return func(state *terraform.State) error {
		is := state.RootModule().Resources[path].Primary
		d := is.Attributes[attr]
		b, err := strconv.ParseBool(d)
		if err != nil {
			return err
		}
		if b != value {
			return fmt.Errorf("at %s expected %t but got %t", path, value, b)
		}
		return nil
	}
}
